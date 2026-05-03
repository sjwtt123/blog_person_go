package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/utils"
)

type commentService struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	articleRepo repository.ArticleRepository,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

func (s *commentService) Create(userID uint, articleID uint, req *request.CreateCommentRequest) (*response.CommentResponse, error) {
	article, err := s.articleRepo.FindByID(articleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "文章不存在")
	}

	comment := &entity.Comment{
		ArticleID:   articleID,
		UserID:      userID,
		ParentID:    req.ParentID,
		ReplyToID:   req.ReplyToID,
		ReplyToName: req.ReplyToName,
		Content:     req.Content,
		Status:      1,
	}

	if comment.ParentID > 0 {
		parent, err := s.commentRepo.FindByID(comment.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "父评论不存在")
		}
		if parent.ArticleID != articleID {
			return nil, bizerrors.New(bizerrors.CodeInvalidParam, "父评论不属于该文章")
		}
	}

	err = s.commentRepo.CreateWithArticleCount(comment, articleID)

	if err != nil {
		return nil, err
	}

	return s.buildCommentResponse(comment)
}

func (s *commentService) GetByArticleID(articleID uint, page, size int) (*response.PageResponse, error) {
	page, size, offset := utils.NormalizeAndOffset(page, size)

	comments, total, err := s.commentRepo.ListByArticleID(articleID, offset, size)
	if err != nil {
		return nil, err
	}

	var result []*response.CommentResponse
	for _, c := range comments {
		item, err := s.buildCommentResponse(c)
		if err != nil {
			continue
		}

		replies, err := s.commentRepo.ListRepliesByParentID(c.ID)
		if err == nil && len(replies) > 0 {
			for _, r := range replies {
				replyResp, err := s.buildCommentResponse(r)
				if err == nil {
					item.Replies = append(item.Replies, replyResp)
				}
			}
		}

		result = append(result, item)
	}

	return response.NewPageResponse(result, total, page, size), nil
}

func (s *commentService) ListAllAdmin(req *request.AdminCommentListRequest) (*response.PageResponse, error) {
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	comments, total, err := s.commentRepo.ListAllAdmin(req.ArticleID, req.Keyword, offset, size)
	if err != nil {
		return nil, err
	}

	var result []*response.CommentResponse
	for _, c := range comments {
		item, err := s.buildCommentResponse(c)
		if err != nil {
			continue
		}
		result = append(result, item)
	}

	return response.NewPageResponse(result, total, page, size), nil
}

func (s *commentService) Update(userID uint, commentID uint, req *request.UpdateCommentRequest) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "评论不存在")
	}

	if comment.UserID != userID {
		return bizerrors.New(bizerrors.CodeForbidden, "无权修改此评论")
	}

	comment.Content = req.Content
	return s.commentRepo.Update(comment)
}

func (s *commentService) Delete(userID uint, role string, commentID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "评论不存在")
	}

	if comment.UserID != userID && role != "admin" {
		return bizerrors.New(bizerrors.CodeForbidden, "无权删除此评论")
	}

	return s.deleteCommentWithTransaction(comment)
}

func (s *commentService) buildCommentResponse(c *entity.Comment) (*response.CommentResponse, error) {
	var userResp *response.CommentUser
	if c.User != nil {
		userResp = &response.CommentUser{
			ID:       c.User.ID,
			Username: c.User.Username,
			Avatar:   c.User.Avatar,
			Role:     c.User.Role,
		}
	}

	var articleResp *response.CommentArticle
	if c.Article != nil {
		articleResp = &response.CommentArticle{
			ID:    c.Article.ID,
			Title: c.Article.Title,
		}
	}

	return &response.CommentResponse{
		ID:          c.ID,
		ArticleID:   c.ArticleID,
		UserID:      c.UserID,
		ParentID:    c.ParentID,
		ReplyToID:   c.ReplyToID,
		ReplyToName: c.ReplyToName,
		Content:     c.Content,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		User:        userResp,
		Article:     articleResp,
	}, nil
}

// AdminDelete 管理员删除评论
func (s *commentService) AdminDelete(commentID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "评论不存在")
	}

	return s.deleteCommentWithTransaction(comment)
}

func (s *commentService) deleteCommentWithTransaction(comment *entity.Comment) error {
	return s.commentRepo.DeleteWithArticleCount(comment.ID, comment.ArticleID)
}
