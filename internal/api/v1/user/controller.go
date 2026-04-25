package user

import (
	"strconv"

	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 用户控制器
type Controller struct {
	userService service.UserService
}

// NewController 创建用户控制器
func NewController(userService service.UserService) *Controller {
	return &Controller{
		userService: userService,
	}
}

// GetAbout 获取关于页信息
// @Summary 获取关于页信息
// @Description 获取博客关于页信息
// @Tags 关于
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/about [get]
func (ctrl *Controller) GetAbout(c *gin.Context) {
	data, err := ctrl.userService.GetAbout()
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

// UpdateAbout 更新关于页信息
// @Summary 更新关于页信息
// @Description 管理员更新关于页信息
// @Tags 关于
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.AboutRequest true "关于页信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/about [put]
func (ctrl *Controller) UpdateAbout(c *gin.Context) {
	var req request.AboutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := ctrl.userService.UpdateAbout(&req, userID); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile [get]
func (ctrl *Controller) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	userResp := ctrl.userService.GetUserResponse(user)
	response.Success(c, userResp)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.UpdateUserRequest true "更新信息"
// @Success 200 {object} response.Response
// @Router /api/v1/user/profile [put]
func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.userService.UpdateUser(userID, &req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/v1/user/password [post]
func (ctrl *Controller) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.userService.ChangePassword(userID, &req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// ListUsers 获取用户列表（分页）
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int true "页码" minimum(1)
// @Param size query int true "每页大小" minimum(1) maximum(100)
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users [get]
func (ctrl *Controller) ListUsers(c *gin.Context) {
	var req request.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pageResp, err := ctrl.userService.ListUsers(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, pageResp)
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 管理员创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.AdminUserCreateRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users [post]
func (ctrl *Controller) CreateUser(c *gin.Context) {
	var req request.AdminUserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userResp, err := ctrl.userService.AdminCreateUser(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "用户创建成功", userResp)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 管理员更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param request body request.AdminUserUpdateRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{id} [put]
func (ctrl *Controller) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req request.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.userService.AdminUpdateUser(uint(id), &req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{id} [delete]
func (ctrl *Controller) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.AdminDeleteUser(uint(id)); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}
