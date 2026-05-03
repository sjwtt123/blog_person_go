package repository

// AuthRepository 认证数据访问接口
type AuthRepository interface {
	// StoreEmailCode 存储邮箱验证码
	StoreEmailCode(email, code string) error
	// GetEmailCode 获取邮箱验证码
	GetEmailCode(email string) (string, error)
	// DeleteEmailCode 删除邮箱验证码
	DeleteEmailCode(email string) error
	// CheckCooldown 检查发送冷却时间
	CheckCooldown(email string) (bool, error)
	// SetCooldown 设置发送冷却时间
	SetCooldown(email string) error
}
