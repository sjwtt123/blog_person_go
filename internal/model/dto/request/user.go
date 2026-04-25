package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
}

// AdminUserCreateRequest 管理员创建用户请求
type AdminUserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
}

// AdminUserUpdateRequest 管理员更新用户请求
type AdminUserUpdateRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Password string `json:"password" binding:"omitempty,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
	Status   *int   `json:"status" binding:"omitempty,oneof=1 2"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	PageRequest
}

// AboutRequest 关于页更新请求
type AboutRequest struct {
	Name    string `json:"name" binding:"omitempty,max=50"`
	Bio     string `json:"bio" binding:"omitempty,max=500"`
	Avatar  string `json:"avatar" binding:"omitempty,max=255"`
	Content string `json:"content" binding:"omitempty"`
}
