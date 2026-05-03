package request

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// SendCodeRequest 发送邮箱验证码请求
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Code            string `json:"code" binding:"required,len=6"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=50"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8,max=50"`
}
