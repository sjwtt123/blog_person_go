package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateWithArticleCount(comment *entity.Comment, articleID uint) error
	FindByID(id uint) (*entity.Comment, error)
	ListByArticleID(articleID uint, offset, limit int) ([]*entity.Comment, int64, error)
	ListAllAdmin(articleID *uint, keyword string, offset, limit int) ([]*entity.Comment, int64, error)
	ListRepliesByParentID(parentID uint) ([]*entity.Comment, error)
	Update(comment *entity.Comment) error
	DeleteWithArticleCount(commentID uint, articleID uint) error
	DeleteByArticleIDWithCount(tx *gorm.DB, articleID uint) error
}
