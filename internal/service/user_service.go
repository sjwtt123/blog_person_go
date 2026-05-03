package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userService) Register(req *request.RegisterRequest) (*response.UserResponse, error) {
	// 验证两次密码是否一致
	if req.Password != req.ConfirmPassword {
		return nil, bizerrors.New(bizerrors.CodeInvalidParam, "两次输入的密码不一致")
	}

	// 检查用户名是否存在
	err := s.IsExistsUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// 检查邮箱是否存在
	err = s.IsExistsEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Status:   1,      // 默认启用
		Role:     "user", // 默认普通用户
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.GetUserResponse(user), nil
}

// GetUserByID 根据 ID 获取用户
func (s *userService) GetUserByID(id uint) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, bizerrors.ErrUserNotFound
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(id uint, req *request.UpdateUserRequest) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if user.Username != req.Username {
		err = s.IsExistsUsername(req.Username)
		if err != nil {
			return err
		}
		user.Username = req.Username
	}

	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(id uint, req *request.ChangePasswordRequest) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if req.NewPassword != req.ConfirmPassword {
		return bizerrors.New(bizerrors.CodeInvalidParam, "两次输入的密码不一致")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// GetUserResponse 获取用户响应
func (s *userService) GetUserResponse(user *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ListUsers 分页获取用户列表
func (s *userService) ListUsers(req *request.UserListRequest) (*response.PageResponse, error) {
	page, size, offset := utils.NormalizeAndOffset(req.Page, req.Size)

	// 查询数据
	users, total, err := s.userRepo.List(offset, size)
	if err != nil {
		return nil, err
	}

	// 转换为响应 DTO
	list := make([]*response.UserResponse, 0, len(users))
	for _, user := range users {
		list = append(list, s.GetUserResponse(user))
	}

	return response.NewPageResponse(list, total, page, size), nil
}

// AdminCreateUser 管理员创建用户
func (s *userService) AdminCreateUser(req *request.AdminUserCreateRequest) (*response.UserResponse, error) {
	// 检查用户名是否存在
	err := s.IsExistsUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// 检查邮箱是否存在
	err = s.IsExistsEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     req.Role,
		Status:   req.Status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.GetUserResponse(user), nil
}

// AdminUpdateUser 管理员更新用户
func (s *userService) AdminUpdateUser(id uint, req *request.AdminUserUpdateRequest) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := s.validateAndUpdateBasicInfo(user, req); err != nil {
		return err
	}
	s.updateUserPassword(user, req.Password)
	s.updateUserMeta(user, req)

	return s.userRepo.Update(user)
}

func (s *userService) validateAndUpdateBasicInfo(user *entity.User, req *request.AdminUserUpdateRequest) error {
	if user.Username != req.Username {
		if err := s.IsExistsUsername(req.Username); err != nil {
			return err
		}
	}
	user.Username = req.Username

	if user.Email != req.Email {
		if err := s.IsExistsEmail(req.Email); err != nil {
			return err
		}
	}
	user.Email = req.Email
	return nil
}

func (s *userService) updateUserPassword(user *entity.User, password string) {
	if password == "" {
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		user.Password = string(hashedPassword)
	}
}

func (s *userService) updateUserMeta(user *entity.User, req *request.AdminUserUpdateRequest) {
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
}

// AdminDeleteUser 管理员删除用户
func (s *userService) AdminDeleteUser(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(user.ID)
}

// GetAbout 获取关于页信息
func (s *userService) GetAbout() (*response.AboutResponse, error) {
	admin, err := s.userRepo.FindAdmin()
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, bizerrors.New(bizerrors.CodeResourceNotFound, "管理员不存在")
	}

	return &response.AboutResponse{
		Name:    admin.Username,
		Avatar:  admin.Avatar,
		Bio:     admin.Bio,
		Content: admin.AboutContent,
	}, nil
}

// UpdateAbout 更新关于页信息（仅管理员）
func (s *userService) UpdateAbout(req *request.AboutRequest, adminID uint) error {
	admin, err := s.GetUserByID(adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "用户不存在")
	}
	if admin.Role != "admin" {
		return bizerrors.New(bizerrors.CodeForbidden, "无权修改关于页信息")
	}

	if req.Name != "" {
		admin.Username = req.Name
	}
	if req.Bio != "" {
		admin.Bio = req.Bio
	}
	if req.Avatar != "" {
		admin.Avatar = req.Avatar
	}
	if req.Content != "" {
		admin.AboutContent = req.Content
	}

	return s.userRepo.Update(admin)
}

func (s *userService) IsExistsUsername(username string) error {
	isExists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if isExists {
		return bizerrors.ErrUsernameAlreadyExists
	}
	return nil

}
func (s *userService) IsExistsEmail(email string) error {
	isExists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return err
	}
	if isExists {
		return bizerrors.ErrEmailAlreadyExists
	}
	return nil
}
