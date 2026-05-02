package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/utils"

	"gorm.io/gorm"
)

// likeService 点赞服务实现
type likeService struct {
	likeRepo    repository.LikeRepository
	articleRepo repository.ArticleRepository
	db          *gorm.DB
}

// NewLikeService 创建点赞服务实例
func NewLikeService(
	likeRepo repository.LikeRepository,
	articleRepo repository.ArticleRepository,
	db *gorm.DB,
) LikeService {
	return &likeService{
		likeRepo:    likeRepo,
		articleRepo: articleRepo,
		db:          db,
	}
}

// LikeArticle 点赞文章
func (s *likeService) LikeArticle(userID uint, articleID uint) (*response.LikeResponse, error) {
	article, err := s.validateArticleExists(articleID)
	if err != nil {
		return nil, err
	}

	// 检查是否已点赞
	existingLike, err := s.likeRepo.FindByUserAndTarget(userID, articleID, "article")
	if err != nil {
		return nil, err
	}
	if existingLike != nil {
		return nil, bizerrors.New(bizerrors.CodeResourceAlreadyExists, "已点赞")
	}

	// 创建点赞记录并更新文章点赞数
	err = s.createLikeAndUpdateCount(userID, articleID, article)
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

	// 检查点赞记录是否存在
	existingLike, err := s.likeRepo.FindByUserAndTarget(userID, articleID, "article")
	if err != nil {
		return nil, err
	}
	if existingLike == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "点赞记录不存在")
	}

	// 删除点赞记录并更新文章点赞数
	err = s.deleteLikeAndUpdateCount(userID, articleID)
	if err != nil {
		return nil, err
	}

	return s.buildLikeResponse(userID, articleID)
}

// GetArticleLikeStatus 获取文章点赞状态
func (s *likeService) GetArticleLikeStatus(userID uint, articleID uint) (*response.LikeResponse, error) {
	_, err := s.validateArticleExists(articleID)
	if err != nil {
		return nil, err
	}

	return s.buildLikeResponse(userID, articleID)
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

// createLikeAndUpdateCount 创建点赞记录并更新文章点赞数（事务操作）
func (s *likeService) createLikeAndUpdateCount(userID uint, articleID uint, article *entity.Article) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		like := &entity.Like{
			UserID:     userID,
			TargetID:   articleID,
			TargetType: "article",
		}
		if err := s.likeRepo.Create(like); err != nil {
			return err
		}

		newCount := article.LikeCount + 1
		return s.articleRepo.UpdateLikeCountInTx(tx, articleID, newCount)
	})
}

// deleteLikeAndUpdateCount 删除点赞记录并更新文章点赞数（事务操作）
func (s *likeService) deleteLikeAndUpdateCount(userID uint, articleID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.likeRepo.Delete(userID, articleID, "article"); err != nil {
			return err
		}

		article, err := s.articleRepo.FindArticleByID(articleID)
		if err != nil {
			return err
		}

		newCount := article.LikeCount - 1
		if newCount < 0 {
			newCount = 0
		}
		return s.articleRepo.UpdateLikeCountInTx(tx, articleID, newCount)
	})
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
