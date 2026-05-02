package like

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller 点赞控制器
type Controller struct {
	likeService service.LikeService
}

// NewController 创建点赞控制器实例
func NewController(likeService service.LikeService) *Controller {
	return &Controller{
		likeService: likeService,
	}
}

// LikeArticle 点赞文章
func (ctrl *Controller) LikeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	data, err := ctrl.likeService.LikeArticle(userID, uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, data)
}

// UnlikeArticle 取消点赞文章
func (ctrl *Controller) UnlikeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	data, err := ctrl.likeService.UnlikeArticle(userID, uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, data)
}

// GetArticleLikeCount 获取文章点赞数（公开接口，无需登录）
func (ctrl *Controller) GetArticleLikeCount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}

	count, err := ctrl.likeService.GetArticleLikeCount(uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, gin.H{"like_count": count})
}

// GetUserLikeStatus 获取用户点赞状态（需要登录）
func (ctrl *Controller) GetUserLikeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	isLiked, err := ctrl.likeService.GetUserLikeStatus(userID, uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, gin.H{"is_liked": isLiked})
}

// ListUserLikedArticles 获取用户点赞的文章列表
func (ctrl *Controller) ListUserLikedArticles(c *gin.Context) {
	var req request.ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	data, err := ctrl.likeService.ListUserLikedArticles(userID, &req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, data)
}
