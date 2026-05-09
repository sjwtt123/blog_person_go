package repository

import (
	"blog/internal/model/entity"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type articleRepository struct {
	db           *gorm.DB
	rdb          *redis.Client
	categoryRepo CategoryRepository
	tagRepo      TagRepository
	commentRepo  CommentRepository
}

func NewArticleRepository(db *gorm.DB, rdb *redis.Client, categoryRepo CategoryRepository, tagRepo TagRepository, commentRepo CommentRepository) ArticleRepository {
	return &articleRepository{
		db:           db,
		rdb:          rdb,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		commentRepo:  commentRepo,
	}
}

func (r *articleRepository) FindByID(id uint) (*entity.Article, error) {
	var a entity.Article
	err := r.db.First(&a, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *articleRepository) FindArticleByID(id uint) (*entity.Article, error) {
	var a entity.Article
	err := r.db.Preload("Author").Preload("Category").Preload("Tags").First(&a, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *articleRepository) List(filter ArticleListFilter) ([]*entity.Article, int64, error) {
	db := r.db.Model(&entity.Article{})

	if filter.TagID != nil && *filter.TagID > 0 {
		db = db.Joins("JOIN article_tags at ON at.article_id = articles.id").
			Where("at.tag_id = ?", *filter.TagID)
	}

	if filter.Keyword != "" {
		db = db.Where("articles.title LIKE ? OR articles.summary LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	if filter.CategoryID != nil && *filter.CategoryID > 0 {
		db = db.Where("articles.category_id = ?", *filter.CategoryID)
	}

	if filter.Status != nil {
		db = db.Where("articles.status = ?", *filter.Status)
	}

	if filter.OnlyPublic {
		db = db.Where("articles.status = 1").Where("articles.published_at IS NOT NULL")
	}

	orderClause := "articles.published_at DESC, articles.created_at DESC"
	if filter.SortBy == "view_count" {
		orderClause = "articles.view_count DESC, articles.published_at DESC"
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*entity.Article
	query := db.Preload("Author").Preload("Category").Preload("Tags").
		Order(orderClause).
		Offset(filter.Offset).
		Limit(filter.Limit)

	if err := query.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// CreateWithCascade 级联创建文章，同时更新分类和标签的文章计数
func (r *articleRepository) CreateWithCascade(article *entity.Article) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// 更新分类文章计数
		if article.CategoryID != nil {
			if err := r.incrementCategoryPostCount(tx, *article.CategoryID); err != nil {
				return err
			}
		}

		// 更新标签文章计数
		for _, tag := range article.Tags {
			if err := r.incrementTagPostCount(tx, tag.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateWithCascade 级联更新文章，处理分类和标签变更时的文章计数
func (r *articleRepository) UpdateWithCascade(article *entity.Article, oldCategoryID *uint, oldTags []*entity.Tag) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 保存文章基本信息（不保存关联）
		if err := tx.Save(article).Error; err != nil {
			return err
		}

		// 替换标签关联（先删除旧关联，再添加新关联）
		if err := tx.Model(article).Association("Tags").Replace(article.Tags); err != nil {
			return err
		}

		// 处理分类变更
		if err := r.handleCategoryChange(tx, oldCategoryID, article.CategoryID); err != nil {
			return err
		}

		// 处理标签变更
		if err := r.handleTagChange(tx, oldTags, article.Tags); err != nil {
			return err
		}

		return nil
	})
}

// DeleteWithCascade 级联删除文章，同时更新分类和标签的文章计数，并删除相关评论
func (r *articleRepository) DeleteWithCascade(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 查询文章详情（用于更新计数）
		article, err := r.FindArticleByID(id)
		if err != nil {
			return err
		}
		if article == nil {
			return nil
		}

		// 删除文章相关评论
		if err := r.commentRepo.DeleteByArticleIDWithCount(tx, id); err != nil {
			return err
		}

		// 删除文章
		if err := tx.Select("Tags").Unscoped().Delete(&entity.Article{}, id).Error; err != nil {
			return err
		}

		// 更新分类文章计数
		if article.CategoryID != nil {
			if err := r.decrementCategoryPostCount(tx, *article.CategoryID); err != nil {
				return err
			}
		}

		// 更新标签文章计数
		for _, tag := range article.Tags {
			if err := r.decrementTagPostCount(tx, tag.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

// incrementCategoryPostCount 增加分类文章计数
func (r *articleRepository) incrementCategoryPostCount(tx *gorm.DB, categoryID uint) error {
	return tx.Model(&entity.Category{}).Where("id = ?", categoryID).Update("post_count", gorm.Expr("post_count + 1")).Error
}

// decrementCategoryPostCount 减少分类文章计数
func (r *articleRepository) decrementCategoryPostCount(tx *gorm.DB, categoryID uint) error {
	return tx.Model(&entity.Category{}).Where("id = ?", categoryID).
		Where("post_count > 0").
		Update("post_count", gorm.Expr("post_count - 1")).Error
}

// incrementTagPostCount 增加标签文章计数
func (r *articleRepository) incrementTagPostCount(tx *gorm.DB, tagID uint) error {
	return tx.Model(&entity.Tag{}).Where("id = ?", tagID).Update("post_count", gorm.Expr("post_count + 1")).Error
}

// decrementTagPostCount 减少标签文章计数
func (r *articleRepository) decrementTagPostCount(tx *gorm.DB, tagID uint) error {
	return tx.Model(&entity.Tag{}).Where("id = ?", tagID).
		Where("post_count > 0").
		Update("post_count", gorm.Expr("post_count - 1")).Error
}

// handleCategoryChange 处理分类变更时的文章计数调整
func (r *articleRepository) handleCategoryChange(tx *gorm.DB, oldCategoryID, newCategoryID *uint) error {
	// 分类发生变化：旧分类减 1，新分类加 1
	if oldCategoryID != nil && newCategoryID != nil && *oldCategoryID != *newCategoryID {
		if err := r.decrementCategoryPostCount(tx, *oldCategoryID); err != nil {
			return err
		}
		return r.incrementCategoryPostCount(tx, *newCategoryID)
	}
	// 新增分类
	if oldCategoryID == nil && newCategoryID != nil {
		return r.incrementCategoryPostCount(tx, *newCategoryID)
	}
	// 移除分类
	if oldCategoryID != nil && newCategoryID == nil {
		return r.decrementCategoryPostCount(tx, *oldCategoryID)
	}
	return nil
}

// handleTagChange 处理标签变更时的文章计数调整
func (r *articleRepository) handleTagChange(tx *gorm.DB, oldTags, newTags []*entity.Tag) error {
	oldTagMap := make(map[uint]bool)
	for _, t := range oldTags {
		oldTagMap[t.ID] = true
	}
	newTagMap := make(map[uint]bool)
	for _, t := range newTags {
		newTagMap[t.ID] = true
	}

	// 移除的标签：文章计数减 1
	for _, t := range oldTags {
		if !newTagMap[t.ID] {
			if err := r.decrementTagPostCount(tx, t.ID); err != nil {
				return err
			}
		}
	}

	// 新增的标签：文章计数加 1
	for _, t := range newTags {
		if !oldTagMap[t.ID] {
			if err := r.incrementTagPostCount(tx, t.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *articleRepository) ListTimeline() ([]*entity.Article, error) {
	var list []*entity.Article
	err := r.db.Model(&entity.Article{}).
		Select("id, title, slug, cover_image, view_count, category_id, created_at").
		Preload("Category").
		Where("status = 1").
		Where("published_at IS NOT NULL").
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *articleRepository) UpdateViewCount(id uint, count uint) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).Update("view_count", count).Error
}

func (r *articleRepository) ListAll() ([]*entity.Article, error) {
	var list []*entity.Article
	err := r.db.Model(&entity.Article{}).
		Select("id", "content", "cover_image").
		Find(&list).Error
	return list, err
}
