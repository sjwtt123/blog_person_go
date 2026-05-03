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
}

func NewArticleRepository(db *gorm.DB, rdb *redis.Client, categoryRepo CategoryRepository, tagRepo TagRepository) ArticleRepository {
	return &articleRepository{
		db:           db,
		rdb:          rdb,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
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

func (r *articleRepository) Create(article *entity.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) CreateWithCascade(article *entity.Article) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}

	if article.CategoryID != nil {
		cat, err := r.categoryRepo.FindByID(*article.CategoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if cat != nil {
			if err := tx.Model(&entity.Category{}).Where("id = ?", *article.CategoryID).Update("post_count", cat.PostCount+1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for _, tag := range article.Tags {
		t, err := r.tagRepo.FindTagByID(tag.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if t != nil {
			if err := tx.Model(&entity.Tag{}).Where("id = ?", tag.ID).Update("post_count", t.PostCount+1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *articleRepository) Update(article *entity.Article) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(article).Error
}

func (r *articleRepository) UpdateWithCascade(article *entity.Article, oldCategoryID *uint, oldTags []*entity.Tag) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(article).Error; err != nil {
		tx.Rollback()
		return err
	}

	if oldCategoryID != nil && article.CategoryID != nil && *oldCategoryID != *article.CategoryID {
		if oldCat, err := r.categoryRepo.FindByID(*oldCategoryID); err == nil && oldCat != nil && oldCat.PostCount > 0 {
			if err := tx.Model(&entity.Category{}).Where("id = ?", *oldCategoryID).Update("post_count", oldCat.PostCount-1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if newCat, err := r.categoryRepo.FindByID(*article.CategoryID); err == nil && newCat != nil {
			if err := tx.Model(&entity.Category{}).Where("id = ?", *article.CategoryID).Update("post_count", newCat.PostCount+1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else if oldCategoryID == nil && article.CategoryID != nil {
		if newCat, err := r.categoryRepo.FindByID(*article.CategoryID); err == nil && newCat != nil {
			if err := tx.Model(&entity.Category{}).Where("id = ?", *article.CategoryID).Update("post_count", newCat.PostCount+1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else if oldCategoryID != nil && article.CategoryID == nil {
		if oldCat, err := r.categoryRepo.FindByID(*oldCategoryID); err == nil && oldCat != nil && oldCat.PostCount > 0 {
			if err := tx.Model(&entity.Category{}).Where("id = ?", *oldCategoryID).Update("post_count", oldCat.PostCount-1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	oldTagMap := make(map[uint]bool)
	for _, t := range oldTags {
		oldTagMap[t.ID] = true
	}
	newTagMap := make(map[uint]bool)
	for _, t := range article.Tags {
		newTagMap[t.ID] = true
	}

	for _, t := range oldTags {
		if !newTagMap[t.ID] {
			if tag, err := r.tagRepo.FindTagByID(t.ID); err == nil && tag != nil && tag.PostCount > 0 {
				if err := tx.Model(&entity.Tag{}).Where("id = ?", t.ID).Update("post_count", tag.PostCount-1).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	for _, t := range article.Tags {
		if !oldTagMap[t.ID] {
			if tag, err := r.tagRepo.FindTagByID(t.ID); err == nil && tag != nil {
				if err := tx.Model(&entity.Tag{}).Where("id = ?", t.ID).Update("post_count", tag.PostCount+1).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	return tx.Commit().Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Select("Tags").Unscoped().Delete(&entity.Article{}, id).Error
}

func (r *articleRepository) DeleteWithCascade(id uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	article, err := r.FindArticleByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if article == nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Select("Tags").Unscoped().Delete(&entity.Article{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if article.CategoryID != nil {
		cat, err := r.categoryRepo.FindByID(*article.CategoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if cat != nil && cat.PostCount > 0 {
			if err := tx.Model(&entity.Category{}).Where("id = ?", *article.CategoryID).Update("post_count", cat.PostCount-1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for _, tag := range article.Tags {
		t, err := r.tagRepo.FindTagByID(tag.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if t != nil && t.PostCount > 0 {
			if err := tx.Model(&entity.Tag{}).Where("id = ?", tag.ID).Update("post_count", t.PostCount-1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
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

func (r *articleRepository) UpdateCommentCount(id uint, count int) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).Update("comment_count", count).Error
}

func (r *articleRepository) UpdateLikeCount(id uint, count int) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).Update("like_count", count).Error
}

func (r *articleRepository) ListAll() ([]*entity.Article, error) {
	var list []*entity.Article
	err := r.db.Model(&entity.Article{}).
		Select("id", "content", "cover_image").
		Find(&list).Error
	return list, err
}
