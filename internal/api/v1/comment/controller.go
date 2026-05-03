package comment

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	commentService service.CommentService
	userService    service.UserService
}

func NewController(commentService service.CommentService, userService service.UserService) *Controller {
	return &Controller{
		commentService: commentService,
		userService:    userService,
	}
}

func (ctrl *Controller) List(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || articleID == 0 {
		response.BadRequest(c, "文章ID参数错误")
		return
	}

	var req request.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := ctrl.commentService.GetByArticleID(uint(articleID), req.Page, req.Size)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) Create(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || articleID == 0 {
		response.BadRequest(c, "文章ID参数错误")
		return
	}

	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	data, err := ctrl.commentService.Create(userID, uint(articleID), &req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) Update(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || commentID == 0 {
		response.BadRequest(c, "评论ID参数错误")
		return
	}

	var req request.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := ctrl.commentService.Update(userID, uint(commentID), &req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || commentID == 0 {
		response.BadRequest(c, "评论ID参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	if err := ctrl.commentService.Delete(userID, role, uint(commentID)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (ctrl *Controller) ListAllAdmin(c *gin.Context) {
	var req request.AdminCommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := ctrl.commentService.ListAllAdmin(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

func (ctrl *Controller) AdminDelete(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || commentID == 0 {
		response.BadRequest(c, "评论ID参数错误")
		return
	}

	if err := ctrl.commentService.AdminDelete(uint(commentID)); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
