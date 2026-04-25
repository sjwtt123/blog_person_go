package repository

import (
	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByID(id uint) (*entity.Category, error)
	FindBySlug(slug string) (*entity.Category, error)
	ListAll() ([]*entity.Category, error)
	Create(category *entity.Category) error
	Update(category *entity.Category) error
	Delete(id uint) error
	UpdatePostCountInTx(tx *gorm.DB, id uint, count int) error
}
