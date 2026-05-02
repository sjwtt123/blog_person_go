package repository

import "blog/internal/model/entity"

type LikeRepository interface {
	Create(like *entity.Like) error
	Delete(userID uint, targetID uint, targetType string) error
	FindByUserAndTarget(userID uint, targetID uint, targetType string) (*entity.Like, error)
	CountByTarget(targetID uint, targetType string) (int64, error)
	ListArticlesByUser(userID uint, offset int, limit int) ([]*entity.Article, int64, error)
}
