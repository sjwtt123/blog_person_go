package response

import "time"

// UserResponse 用户响应
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	Bio       string     `json:"bio"`
	Role      string     `json:"role"`
	Status    int        `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// AboutResponse 关于页响应
type AboutResponse struct {
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Bio     string `json:"bio"`
	Content string `json:"content"`
}

// AboutRequest 关于页更新请求
type AboutRequest struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Avatar  string `json:"avatar"`
	Content string `json:"content"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}
