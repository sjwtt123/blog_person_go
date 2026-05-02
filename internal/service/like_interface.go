package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
)

type LikeService interface {
	LikeArticle(userID uint, articleID uint) (*response.LikeResponse, error)
	UnlikeArticle(userID uint, articleID uint) (*response.LikeResponse, error)
	GetArticleLikeStatus(userID uint, articleID uint) (*response.LikeResponse, error)
	ListUserLikedArticles(userID uint, req *request.ArticleListRequest) (*response.PageResponse, error)
}
