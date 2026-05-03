package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
)

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(req *request.LoginRequest) (*dto.LoginResponse, error)
	// RefreshToken 刷新 Token
	RefreshToken(token string) (string, error)
	// SendEmailCode 发送邮箱验证码
	SendEmailCode(email string) error
	// ResetPassword 重置密码
	ResetPassword(req *request.ResetPasswordRequest) error
	// VerifyEmailCode 验证邮箱验证码
	VerifyEmailCode(email, code string) error
}
