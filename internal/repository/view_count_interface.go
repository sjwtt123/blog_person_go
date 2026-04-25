package repository

import (
	"time"
)

// ViewCountRepository 文章访问量统计仓库接口
type ViewCountRepository interface {
	Increment(articleID uint) error
	Get(articleID uint) (int64, error)
	Delete(articleID uint) error
	GetAll() (map[string]string, error)
	SetExpire(duration time.Duration) error
}
