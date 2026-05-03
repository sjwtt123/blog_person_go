package service

import (
	"blog/pkg/captcha"
	"blog/pkg/logger"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"regexp"

	"go.uber.org/zap"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/config"
	bizerrors "blog/pkg/errors"
	"blog/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// 密码强度正则表达式（必须包含数字和字母）
var (
	passwordDigitRegex  = regexp.MustCompile(`[0-9]`)
	passwordLetterRegex = regexp.MustCompile(`[a-zA-Z]`)
)

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository // 用户数据访问层
	authRepo repository.AuthRepository // 认证数据访问层（Redis验证码操作）
	userSvc  UserService               // 用户服务
	cfg      config.EmailConfig        // 邮件配置
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, authRepo repository.AuthRepository, userSvc UserService, cfg config.EmailConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		authRepo: authRepo,
		userSvc:  userSvc,
		cfg:      cfg,
	}
}

// Login 用户登录
func (s *authService) Login(req *request.LoginRequest) (*dto.LoginResponse, error) {
	// 验证图形验证码
	if !captcha.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		return nil, bizerrors.New(bizerrors.CodeInvalidParam, "验证码错误")
	}

	// 根据用户名查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, bizerrors.ErrInvalidCredentials
	}

	// 检查用户状态是否启用
	if err := s.checkUserStatus(user.Status); err != nil {
		return nil, err
	}

	// 验证密码是否正确
	if err := s.verifyPassword(user.Password, req.Password); err != nil {
		return nil, bizerrors.ErrInvalidCredentials
	}

	// 生成JWT Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 返回用户信息和Token
	return s.buildLoginResponse(user, token)
}

// RefreshToken 刷新Token
func (s *authService) RefreshToken(token string) (string, error) {
	return jwt.RefreshToken(token)
}

// SendEmailCode 发送邮箱验证码
func (s *authService) SendEmailCode(email string) error {
	// 检查发送冷却时间（60秒内不能重复发送）
	if err := s.checkCooldown(email); err != nil {
		return err
	}

	// 生成6位数字验证码
	code, err := s.generateCode()
	if err != nil {
		return err
	}

	// 将验证码存储到Redis（5分钟过期）
	if err := s.authRepo.StoreEmailCode(email, code); err != nil {
		return err
	}

	// 设置发送冷却时间
	if err := s.authRepo.SetCooldown(email); err != nil {
		return err
	}

	// 发送邮件
	return s.sendEmail(email, code)
}

// VerifyEmailCode 验证邮箱验证码
func (s *authService) VerifyEmailCode(email, code string) error {
	// 从Redis获取存储的验证码
	storedCode, err := s.authRepo.GetEmailCode(email)
	if errors.Is(err, redis.Nil) {
		return bizerrors.New(bizerrors.CodeInvalidParam, "验证码已过期")
	}
	if err != nil {
		return err
	}
	// 比对验证码是否正确
	if storedCode != code {
		return bizerrors.New(bizerrors.CodeInvalidParam, "验证码错误")
	}
	return nil
}

// ResetPassword 重置密码
func (s *authService) ResetPassword(req *request.ResetPasswordRequest) error {
	// 验证两次密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		return bizerrors.New(bizerrors.CodeInvalidParam, "两次输入的密码不一致")
	}

	// 验证新密码强度（包含数字和字母）
	if err := s.validatePassword(req.NewPassword); err != nil {
		return err
	}

	// 验证邮箱验证码
	if err := s.VerifyEmailCode(req.Email, req.Code); err != nil {
		return err
	}

	// 根据邮箱查找用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return bizerrors.New(bizerrors.CodeResourceNotFound, "用户不存在")
	}

	// 对新密码进行bcrypt加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新用户密码
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 删除Redis中的验证码
	return s.authRepo.DeleteEmailCode(req.Email)
}

// checkUserStatus 检查用户状态
func (s *authService) checkUserStatus(status int) error {
	if status != 1 {
		return bizerrors.ErrUserDisabled
	}
	return nil
}

// verifyPassword 验证密码（使用bcrypt比较）
func (s *authService) verifyPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// buildLoginResponse 构建登录响应（包含Token和用户信息）
func (s *authService) buildLoginResponse(user *entity.User, token string) (*dto.LoginResponse, error) {
	userResp := s.userSvc.GetUserResponse(user)
	return &dto.LoginResponse{
		Token: token,
		User:  *userResp,
	}, nil
}

// checkCooldown 检查发送冷却时间（防止频繁发送）
func (s *authService) checkCooldown(email string) error {
	inCooldown, err := s.authRepo.CheckCooldown(email)
	if err != nil {
		return err
	}
	if inCooldown {
		return bizerrors.New(bizerrors.CodeTooManyRequests, "请60秒后再发送")
	}
	return nil
}

// generateCode 生成6位数字验证码（使用加密安全的随机数）
func (s *authService) generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", bizerrors.New(bizerrors.CodeInternalError, "验证码生成失败")
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// validatePassword 验证密码强度（要求包含数字和字母）
func (s *authService) validatePassword(password string) error {
	if len(password) == 0 || !passwordDigitRegex.MatchString(password) || !passwordLetterRegex.MatchString(password) {
		return bizerrors.New(bizerrors.CodeInvalidParam, "密码必须包含数字和字母")
	}
	return nil
}

// sendEmail 发送邮件（根据端口自动选择连接方式）
func (s *authService) sendEmail(to string, code string) error {
	msg := s.buildEmailMessage(to, code)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	// 端口465使用SSL连接
	if s.cfg.Port == 465 {
		return s.sendEmailSSL(addr, to, msg)
	}

	// 其他端口使用STARTTLS
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	return smtp.SendMail(addr, auth, s.cfg.Username, []string{to}, []byte(msg))
}

// sendEmailSSL 通过SSL发送邮件（端口465）
func (s *authService) sendEmailSSL(addr, to string, msg string) error {
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	// 建立TLS连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.cfg.Host})
	if err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "连接邮件服务器失败")
	}
	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error("关闭连接邮件服务器失败", zap.Error(err))
		}
	}(conn)

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "创建邮件客户端失败")
	}
	defer func(client *smtp.Client) {
		err := client.Close()
		if err != nil {
			logger.Error("关闭连接邮件客户端失败", zap.Error(err))
		}
	}(client)

	// 身份认证
	if err := client.Auth(auth); err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "邮件认证失败")
	}

	// 设置发件人
	if err := client.Mail(s.cfg.Username); err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "设置发件人失败")
	}

	// 设置收件人
	if err := client.Rcpt(to); err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "设置收件人失败")
	}

	// 获取邮件数据流
	w, err := client.Data()
	if err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "获取邮件数据流失败")
	}

	// 写入邮件内容
	if _, err := w.Write([]byte(msg)); err != nil {
		return bizerrors.New(bizerrors.CodeInternalError, "写入邮件内容失败")
	}

	// 关闭数据流并发送
	return w.Close()
}

// buildEmailMessage 构建符合RFC 822标准的邮件内容
func (s *authService) buildEmailMessage(to string, code string) string {
	subject := "【Blog】邮箱验证码"
	body := fmt.Sprintf("您的验证码是：%s，有效期为5分钟，请勿泄露给他人。", code)
	return fmt.Sprintf("From: \"%s\" <%s>\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, s.cfg.Username, to, subject, body)
}
