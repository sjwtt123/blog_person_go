package repository

import (
	"blog/internal/model/entity"
	"errors"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {

	var c entity.Category

	err := r.db.First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) FindBySlug(slug string) (*entity.Category, error) {
	var c entity.Category
	err := r.db.Where("slug = ?", slug).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) ListAll() ([]*entity.Category, error) {
	var list []*entity.Category
	if err := r.db.Order("sort_order ASC, created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}

func (r *categoryRepository) UpdatePostCountInTx(tx *gorm.DB, id uint, count int) error {
	return tx.Model(&entity.Category{}).Where("id = ?", id).Update("post_count", count).Error
}
