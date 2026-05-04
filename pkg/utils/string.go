package utils

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
)

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// Contains 检查字符串是否在切片中
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

var (
	nonSlugChar = regexp.MustCompile(`[^a-z0-9\-]+`)
	multiDash   = regexp.MustCompile(`-+`)
)

// Slugify 将字符串转换为 slug 格式
func Slugify(s string) string {
	v := strings.ToLower(strings.TrimSpace(s))
	v = strings.ReplaceAll(v, " ", "-")
	v = nonSlugChar.ReplaceAllString(v, "-")
	v = multiDash.ReplaceAllString(v, "-")
	v = strings.Trim(v, "-")
	return v
}
