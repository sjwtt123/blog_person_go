package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"

	"blog/internal/model/entity"
)

// UserService 用户服务接口
type UserService interface {

	// Register 用户注册
	Register(req *request.RegisterRequest) (*response.UserResponse, error)

	// GetUserByID 根据 ID 获取用户
	GetUserByID(id uint) (*entity.User, error)

	// UpdateUser 更新用户信息
	UpdateUser(id uint, req *request.UpdateUserRequest) error

	// ChangePassword 修改密码
	ChangePassword(id uint, req *request.ChangePasswordRequest) error

	// GetUserResponse 获取用户响应
	GetUserResponse(user *entity.User) *response.UserResponse

	// ListUsers 分页获取用户列表
	ListUsers(req *request.UserListRequest) (*response.PageResponse, error)

	// AdminCreateUser 管理员创建用户
	AdminCreateUser(req *request.AdminUserCreateRequest) (*response.UserResponse, error)

	// AdminUpdateUser 管理员更新用户
	AdminUpdateUser(id uint, req *request.AdminUserUpdateRequest) error

	// AdminDeleteUser 管理员删除用户
	AdminDeleteUser(id uint) error

	// GetAbout 获取关于页信息
	GetAbout() (*response.AboutResponse, error)

	// UpdateAbout 更新关于页信息（仅管理员）
	UpdateAbout(req *request.AboutRequest, adminID uint) error
}
