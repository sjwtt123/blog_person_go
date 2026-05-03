package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/utils"
)

type likeService struct {
	likeRepo    repository.LikeRepository
	articleRepo repository.ArticleRepository
}

func NewLikeService(
	likeRepo repository.LikeRepository,
	articleRepo repository.ArticleRepository,
) LikeService {
	return &likeService{
		likeRepo:    likeRepo,
		articleRepo: articleRepo,
	}
}

// LikeArticle 点赞文章
func (s *likeService) LikeArticle(userID uint, articleID uint) (*response.LikeResponse, error) {
	_, err := s.validateArticleExists(articleID)
	if err != nil {
		return nil, err
	}

	existingLike, err := s.likeRepo.FindByUserAndTarget(userID, articleID, "article")
	if err != nil {
		return nil, err
	}
	if existingLike != nil {
		return nil, bizerrors.New(bizerrors.CodeResourceAlreadyExists, "已点赞")
	}

	err = s.likeRepo.CreateLikeAndUpdateCount(userID, articleID)
	if err != nil {
		return nil, err
	}

	return s.buildLikeResponse(userID, articleID)
}

// UnlikeArticle 取消点赞文章
func (s *likeService) UnlikeArticle(userID uint, articleID uint) (*response.LikeResponse, error) {
	_, err := s.validateArticleExists(articleID)
	if err != nil {
		return nil, err
	}

	existingLike, err := s.likeRepo.FindByUserAndTarget(userID, articleID, "article")
	if err != nil {
		return nil, err
	}
	if existingLike == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "点赞记录不存在")
	}

	err = s.likeRepo.DeleteLikeAndUpdateCount(userID, articleID)
	if err != nil {
		return nil, err
	}

	return s.buildLikeResponse(userID, articleID)
}

// GetArticleLikeCount 获取文章点赞数（公开接口）
func (s *likeService) GetArticleLikeCount(articleID uint) (int, error) {
	article, err := s.validateArticleExists(articleID)
	if err != nil {
		return 0, err
	}
	return article.LikeCount, nil
}

// GetUserLikeStatus 获取用户点赞状态（需登录）
func (s *likeService) GetUserLikeStatus(userID uint, articleID uint) (bool, error) {
	_, err := s.validateArticleExists(articleID)
	if err != nil {
		return false, err
	}

	like, err := s.likeRepo.FindByUserAndTarget(userID, articleID, "article")
	if err != nil {
		return false, err
	}
	return like != nil, nil
}

// ListUserLikedArticles 获取用户点赞的文章列表
func (s *likeService) ListUserLikedArticles(userID uint, req *request.ArticleListRequest) (*response.PageResponse, error) {
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	articles, total, err := s.likeRepo.ListArticlesByUser(userID, offset, size)
	if err != nil {
		return nil, err
	}

	var result []*response.ArticleListItemResponse
	for _, article := range articles {
		item := toArticleListItem(article)
		result = append(result, &item)
	}

	return response.NewPageResponse(result, total, page, size), nil
}

// validateArticleExists 验证文章是否存在
func (s *likeService) validateArticleExists(articleID uint) (*entity.Article, error) {
	article, err := s.articleRepo.FindArticleByID(articleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}
	return article, nil
}

// buildLikeResponse 构建点赞响应数据
func (s *likeService) buildLikeResponse(userID uint, articleID uint) (*response.LikeResponse, error) {
	count, err := s.likeRepo.CountByTarget(articleID, "article")
	if err != nil {
		return nil, err
	}

	isLiked := false
	if userID > 0 {
		like, err := s.likeRepo.FindByUserAndTarget(userID, articleID, "article")
		if err != nil {
			return nil, err
		}
		isLiked = like != nil
	}

	return &response.LikeResponse{
		LikeCount: int(count),
		IsLiked:   isLiked,
	}, nil
}
