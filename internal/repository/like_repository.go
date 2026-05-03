package repository

import (
	"blog/internal/model/entity"
	"errors"

	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *entity.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(userID uint, targetID uint, targetType string) error {
	return r.db.Model(&entity.Like{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&entity.Like{}).Error
}

func (r *likeRepository) FindByUserAndTarget(userID uint, targetID uint, targetType string) (*entity.Like, error) {
	var like entity.Like
	err := r.db.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) CountByTarget(targetID uint, targetType string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Count(&count).Error
	return count, err
}

func (r *likeRepository) ListArticlesByUser(userID uint, offset int, limit int) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	query := r.db.Model(&entity.Article{}).
		Joins("INNER JOIN likes ON likes.target_id = articles.id").
		Where("likes.user_id = ? AND likes.target_type = ?", userID, "article").
		Where("articles.status = 1")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("likes.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *likeRepository) CreateLikeAndUpdateCount(userID uint, articleID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	like := &entity.Like{
		UserID:     userID,
		TargetID:   articleID,
		TargetType: "article",
	}
	if err := tx.Create(like).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&entity.Article{}).Where("id = ?", articleID).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *likeRepository) DeleteLikeAndUpdateCount(userID uint, articleID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&entity.Like{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, articleID, "article").
		Delete(&entity.Like{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&entity.Article{}).Where("id = ?", articleID).Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
