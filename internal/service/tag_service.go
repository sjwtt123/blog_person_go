package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/utils"
)

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) ListAll() ([]dto.TagResponse, error) {
	list, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	out := make([]dto.TagResponse, 0, len(list))
	for _, t := range list {
		out = append(out, dto.TagResponse{
			ID:          t.ID,
			Name:        t.Name,
			Slug:        t.Slug,
			Description: t.Description,
			PostCount:   t.PostCount,
		})
	}
	return out, nil
}

func (s *tagService) Create(req *request.TagUpsertRequest) error {
	existing, err := s.repo.FindTagByName(req.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "标签已存在")
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.Slugify(req.Name)
	}

	if slug != "" {
		bySlug, err := s.repo.FindBySlug(slug)
		if err != nil {
			return err
		}
		if bySlug != nil {
			return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "标签 slug 已存在")
		}
	}

	return s.repo.Create(&entity.Tag{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
	})
}

func (s *tagService) Update(id uint, req *request.TagUpsertRequest) error {
	tag, err := s.repo.FindTagByID(id)
	if err != nil {
		return err
	}
	if tag == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "标签不存在")
	}

	if req.Name != tag.Name {
		existing, err := s.repo.FindTagByName(req.Name)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != tag.ID {
			return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "标签名已存在")
		}
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.Slugify(req.Name)
	}
	if slug != "" && slug != tag.Slug {
		bySlug, err := s.repo.FindBySlug(slug)
		if err != nil {
			return err
		}
		if bySlug != nil && bySlug.ID != tag.ID {
			return bizerrors.New(bizerrors.CodeResourceAlreadyExists, "标签 slug 已存在")
		}
	}

	tag.Name = req.Name
	tag.Slug = slug
	tag.Description = req.Description
	return s.repo.Update(tag)
}

func (s *tagService) Delete(id uint) error {
	tag, err := s.repo.FindTagByID(id)
	if err != nil {
		return err
	}
	if tag == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "标签不存在")
	}
	return s.repo.Delete(id)
}
