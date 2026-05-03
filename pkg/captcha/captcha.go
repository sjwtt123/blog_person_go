package captcha

import (
	"time"

	"blog/pkg/database"

	"github.com/mojocn/base64Captcha"
	"golang.org/x/net/context"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成图形验证码
func GenerateCaptcha() (id, b64s string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	return cp.Generate()
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, answer string) bool {
	return store.Verify(id, answer, true)
}

// StoreCaptcha 将验证码存储到 Redis（用于登录防刷）
func StoreCaptcha(ctx context.Context, key, value string, expiration time.Duration) error {
	return database.GetRedis().Set(ctx, key, value, expiration).Err()
}

// GetCaptcha 从 Redis 获取验证码
func GetCaptcha(ctx context.Context, key string) (string, error) {
	return database.GetRedis().Get(ctx, key).Result()
}

// DeleteCaptcha 删除验证码
func DeleteCaptcha(ctx context.Context, key string) error {
	return database.GetRedis().Del(ctx, key).Err()
}
