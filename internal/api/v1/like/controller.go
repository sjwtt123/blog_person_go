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

// GetArticleLikeStatus 获取文章点赞状态（未登录用户仅返回点赞数）
func (ctrl *Controller) GetArticleLikeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 参数错误")
		return
	}

	var userID uint
	if middleware.IsLoggedIn(c) {
		userID = middleware.GetUserID(c)
	}

	data, err := ctrl.likeService.GetArticleLikeStatus(userID, uint(id))
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, data)
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
