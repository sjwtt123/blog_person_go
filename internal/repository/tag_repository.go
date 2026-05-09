package repository

import (
	"errors"

	"blog/internal/model/entity"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindTagByID(id uint) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindTagByName(name string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindBySlug(slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) ListAll() ([]*entity.Tag, error) {
	var tags []*entity.Tag
	if err := r.db.Order("post_count DESC, created_at DESC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) Create(tag *entity.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *entity.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Tag{}, id).Error
}
