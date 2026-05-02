package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
)

type LikeService interface {
	LikeArticle(userID uint, articleID uint) (*response.LikeResponse, error)
	UnlikeArticle(userID uint, articleID uint) (*response.LikeResponse, error)
	GetArticleLikeCount(articleID uint) (int, error)
	GetUserLikeStatus(userID uint, articleID uint) (bool, error)
	ListUserLikedArticles(userID uint, req *request.ArticleListRequest) (*response.PageResponse, error)
}
