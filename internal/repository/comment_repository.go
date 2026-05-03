package repository

import (
	"blog/internal/model/entity"
	"errors"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) CreateWithArticleCount(comment *entity.Comment, articleID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&entity.Article{}).Where("id = ?", articleID).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *commentRepository) FindByID(id uint) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) ListByArticleID(articleID uint, offset, limit int) ([]*entity.Comment, int64, error) {
	var list []*entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).
		Where("article_id = ? AND status = 1", articleID).
		Where("parent_id = 0")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *commentRepository) ListAllAdmin(articleID *uint, keyword string, offset, limit int) ([]*entity.Comment, int64, error) {
	var list []*entity.Comment
	var total int64

	query := r.db.Model(&entity.Comment{}).Where("status = 1")

	if articleID != nil {
		query = query.Where("article_id = ?", *articleID)
	}
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").
		Preload("Article").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *commentRepository) ListRepliesByParentID(parentID uint) ([]*entity.Comment, error) {
	var list []*entity.Comment
	err := r.db.Preload("User").
		Where("parent_id = ? AND status = 1", parentID).
		Order("created_at ASC").
		Find(&list).Error
	return list, err
}

func (r *commentRepository) Update(comment *entity.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Comment{}, id).Error
}

func (r *commentRepository) DeleteWithArticleCount(commentID uint, articleID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&entity.Comment{}, commentID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&entity.Article{}).Where("id = ?", articleID).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *commentRepository) CountByArticleID(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Comment{}).
		Where("article_id = ? AND status = 1", articleID).
		Count(&count).Error
	return count, err
}

func (r *commentRepository) WithTx(tx *gorm.DB) CommentRepository {
	if tx == nil {
		return r
	}
	return &commentRepository{db: tx}
}
