package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type TagRepository interface {
	FindTagByID(id uint) (*entity.Tag, error)
	FindTagByName(name string) (*entity.Tag, error)
	FindBySlug(slug string) (*entity.Tag, error)
	ListAll() ([]*entity.Tag, error)
	Create(tag *entity.Tag) error
	Update(tag *entity.Tag) error
	Delete(id uint) error
	UpdatePostCountInTx(tx *gorm.DB, id uint, count int) error
}
