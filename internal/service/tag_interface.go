package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

type TagService interface {
	ListAll() ([]dto.TagResponse, error)
	Create(req *request.TagUpsertRequest) error
	Update(id uint, req *request.TagUpsertRequest) error
	Delete(id uint) error
}
