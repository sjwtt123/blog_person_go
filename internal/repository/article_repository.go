package repository

import (
	"blog/internal/model/entity"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type articleRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewArticleRepository(db *gorm.DB, rdb *redis.Client) ArticleRepository {
	return &articleRepository{
		db:  db,
		rdb: rdb,
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
		db = db.Where("articles.title LIKE ?", "%"+filter.Keyword+"%")
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

func (r *articleRepository) CreateInTx(tx *gorm.DB, article *entity.Article) error {
	return tx.Create(article).Error
}

func (r *articleRepository) Update(article *entity.Article) error {

	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(article).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Select("Tags").Unscoped().Delete(&entity.Article{}, id).Error
}

func (r *articleRepository) DeleteInTx(tx *gorm.DB, id uint) error {
	return tx.Select("Tags").Unscoped().Delete(&entity.Article{}, id).Error
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

func (r *articleRepository) UpdateViewCountInTx(tx *gorm.DB, id uint, count uint) error {
	return tx.Model(&entity.Article{}).Where("id = ?", id).Update("view_count", count).Error
}

func (r *articleRepository) UpdateCommentCount(id uint, count int) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).Update("comment_count", count).Error
}

func (r *articleRepository) UpdateCommentCountInTx(tx *gorm.DB, id uint, count int) error {
	return tx.Model(&entity.Article{}).Where("id = ?", id).Update("comment_count", count).Error
}

func (r *articleRepository) UpdateLikeCount(id uint, count int) error {
	return r.db.Model(&entity.Article{}).Where("id = ?", id).Update("like_count", count).Error
}

func (r *articleRepository) UpdateLikeCountInTx(tx *gorm.DB, id uint, count int) error {
	return tx.Model(&entity.Article{}).Where("id = ?", id).Update("like_count", count).Error
}

func (r *articleRepository) ListAll() ([]*entity.Article, error) {
	var list []*entity.Article
	err := r.db.Model(&entity.Article{}).
		Select("id", "content", "cover_image").
		Find(&list).Error
	return list, err
}
