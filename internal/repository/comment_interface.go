package repository

import (
	"blog/internal/model/entity"
)

type CommentRepository interface {
	Create(comment *entity.Comment) error
	CreateWithArticleCount(comment *entity.Comment, articleID uint) error
	FindByID(id uint) (*entity.Comment, error)
	ListByArticleID(articleID uint, offset, limit int) ([]*entity.Comment, int64, error)
	ListAllAdmin(articleID *uint, keyword string, offset, limit int) ([]*entity.Comment, int64, error)
	ListRepliesByParentID(parentID uint) ([]*entity.Comment, error)
	Update(comment *entity.Comment) error
	Delete(id uint) error
	DeleteWithArticleCount(commentID uint, articleID uint) error
	CountByArticleID(articleID uint) (int64, error)
}
