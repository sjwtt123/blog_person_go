package repository

import (
	"blog/internal/model/entity"
)

type ArticleListFilter struct {
	Offset     int
	Limit      int
	Keyword    string
	CategoryID *uint
	TagID      *uint
	Status     *int
	OnlyPublic bool
	SortBy     string
}

type ArticleRepository interface {
	FindByID(id uint) (*entity.Article, error)
	FindArticleByID(id uint) (*entity.Article, error)
	List(filter ArticleListFilter) ([]*entity.Article, int64, error)
	ListTimeline() ([]*entity.Article, error)
	Create(article *entity.Article) error
	CreateWithCascade(article *entity.Article) error
	Update(article *entity.Article) error
	UpdateWithCascade(article *entity.Article, oldCategoryID *uint, oldTags []*entity.Tag) error
	Delete(id uint) error
	DeleteWithCascade(id uint) error

	UpdateViewCount(id uint, count uint) error
	UpdateCommentCount(id uint, count int) error
	UpdateLikeCount(id uint, count int) error

	ListAll() ([]*entity.Article, error)
}
