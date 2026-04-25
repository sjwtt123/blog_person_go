package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

type CategoryService interface {
	ListAll() ([]dto.CategoryResponse, error)
	Create(req *request.CategoryUpsertRequest) error
	Update(id uint, req *request.CategoryUpsertRequest) error
	Delete(id uint) error
}
