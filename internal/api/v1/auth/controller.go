package auth

import (
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/captcha"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// Controller 认证控制器
type Controller struct {
	authService service.AuthService
	userService service.UserService
}

// NewController 创建认证控制器
func NewController(authService service.AuthService, userService service.UserService) *Controller {
	return &Controller{
		authService: authService,
		userService: userService,
	}
}

func (ctrl *Controller) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.authService.VerifyEmailCode(req.Email, req.Code); err != nil {
		response.BizError(c, err)
		return
	}

	userResp, err := ctrl.userService.Register(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "注册成功", userResp)
}

// SendCode 发送邮箱验证码
func (ctrl *Controller) SendCode(c *gin.Context) {
	var req request.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.authService.SendEmailCode(req.Email); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "验证码已发送", nil)
}

// ResetPassword 重置密码
func (ctrl *Controller) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.authService.ResetPassword(&req); err != nil {
		response.BizError(c, err)
		return
	}

	response.SuccessWithMessage(c, "密码重置成功", nil)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取 Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=response.LoginResponse}
// @Router /api/v1/auth/login [post]
func (ctrl *Controller) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := ctrl.authService.Login(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetCaptcha 获取图形验证码
func (ctrl *Controller) GetCaptcha(c *gin.Context) {
	id, b64s, err := captcha.GenerateCaptcha()
	if err != nil {
		response.InternalError(c, "验证码生成失败")
		return
	}

	response.Success(c, map[string]string{
		"id":      id,
		"captcha": b64s,
	})
}

// RefreshToken 刷新 Token
// @Summary 刷新 Token
// @Description 使用旧 Token 获取新 Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Token"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (ctrl *Controller) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	newToken, err := ctrl.authService.RefreshToken(req.Token)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, map[string]string{
		"token": newToken,
	})
}
