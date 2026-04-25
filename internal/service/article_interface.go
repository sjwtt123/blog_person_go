package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
)

type ArticleService interface {
	ListPublic(req *request.ArticleListRequest) (*response.PageResponse, error)
	GetDetail(id uint) (*response.ArticleDetailResponse, error)
	GetTimelineItems() ([]*response.TimelineItem, error)

	ListAdmin(req *request.AdminArticleListRequest) (*response.PageResponse, error)
	Create(authorID uint, req *request.ArticleUpsertRequest) error
	Update(id uint, req *request.ArticleUpsertRequest) error
	Delete(id uint) error
}
