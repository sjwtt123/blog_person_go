package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
)

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListAll() ([]dto.CategoryResponse, error) {
	list, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	out := make([]dto.CategoryResponse, 0, len(list))
	for _, c := range list {
		out = append(out, dto.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			PostCount:   c.PostCount,
			SortOrder:   c.SortOrder,
			ParentID:    c.ParentID,
		})
	}
	return out, nil
}

func (s *categoryService) Create(req *request.CategoryUpsertRequest) error {
	existing, err := s.repo.FindBySlug(req.Slug)
	if err != nil {
		return err
	}
	if existing != nil {
		return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "分类已存在")
	}
	return s.repo.Create(&entity.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		ParentID:    req.ParentID,
	})
}

func (s *categoryService) Update(id uint, req *request.CategoryUpsertRequest) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if cat == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "分类不存在")
	}
	if req.Slug != cat.Slug {
		existing, err := s.repo.FindBySlug(req.Slug)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != cat.ID {
			return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "分类 slug 已存在")
		}
	}

	cat.Name = req.Name
	cat.Slug = req.Slug
	cat.Description = req.Description
	cat.SortOrder = req.SortOrder
	cat.ParentID = req.ParentID
	return s.repo.Update(cat)
}

func (s *categoryService) Delete(id uint) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if cat == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "分类不存在")
	}
	return s.repo.Delete(id)
}
