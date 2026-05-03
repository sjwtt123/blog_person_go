package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
)

type CommentService interface {
	Create(userID uint, articleID uint, req *request.CreateCommentRequest) (*response.CommentResponse, error)
	GetByArticleID(articleID uint, page, size int) (*response.PageResponse, error)
	ListAllAdmin(req *request.AdminCommentListRequest) (*response.PageResponse, error)
	Update(userID uint, commentID uint, req *request.UpdateCommentRequest) error
	Delete(userID uint, role string, commentID uint) error
	AdminDelete(commentID uint) error
}
