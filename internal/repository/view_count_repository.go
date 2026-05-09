package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const viewCountHashKey = "article:view_counts"

type viewCountRepository struct {
	rdb *redis.Client
}

func NewViewCountRepository(rdb *redis.Client) ViewCountRepository {
	return &viewCountRepository{
		rdb: rdb,
	}
}

func (r *viewCountRepository) Increment(articleID uint) error {
	ctx := context.Background()
	err := r.rdb.HIncrBy(ctx, viewCountHashKey, strconv.FormatUint(uint64(articleID), 10), 1).Err()
	if err != nil {
		return err
	}

	return r.SetExpire(30 * 24 * time.Hour)
}

func (r *viewCountRepository) Get(articleID uint) (int64, error) {
	ctx := context.Background()
	val, err := r.rdb.HGet(ctx, viewCountHashKey, strconv.FormatUint(uint64(articleID), 10)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}

	return strconv.ParseInt(val, 10, 64)
}

func (r *viewCountRepository) Delete(articleID uint) error {
	ctx := context.Background()
	err := r.rdb.HDel(ctx, viewCountHashKey, strconv.FormatUint(uint64(articleID), 10)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *viewCountRepository) GetAll() (map[string]string, error) {
	ctx := context.Background()
	result, err := r.rdb.HGetAll(ctx, viewCountHashKey).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *viewCountRepository) SetExpire(duration time.Duration) error {
	ctx := context.Background()
	err := r.rdb.Expire(ctx, viewCountHashKey, duration).Err()
	if err != nil {
		return err
	}
	return nil
}
