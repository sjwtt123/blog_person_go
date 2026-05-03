package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	emailCodePrefix   = "email:code:"
	emailCodeExpire   = 5 * time.Minute
	emailCodeCooldown = 60 * time.Second
)

type authRepo struct {
	rdb *redis.Client
}

func NewAuthRepository(rdb *redis.Client) AuthRepository {
	return &authRepo{rdb: rdb}
}

func (r *authRepo) StoreEmailCode(email, code string) error {
	ctx := context.Background()
	key := emailCodePrefix + email
	return r.rdb.Set(ctx, key, code, emailCodeExpire).Err()
}

func (r *authRepo) GetEmailCode(email string) (string, error) {
	ctx := context.Background()
	key := emailCodePrefix + email
	return r.rdb.Get(ctx, key).Result()
}

func (r *authRepo) DeleteEmailCode(email string) error {
	ctx := context.Background()
	key := emailCodePrefix + email
	return r.rdb.Del(ctx, key).Err()
}

func (r *authRepo) CheckCooldown(email string) (bool, error) {
	ctx := context.Background()
	cooldownKey := emailCodePrefix + email + ":cooldown"
	exists, err := r.rdb.Exists(ctx, cooldownKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *authRepo) SetCooldown(email string) error {
	ctx := context.Background()
	cooldownKey := emailCodePrefix + email + ":cooldown"
	return r.rdb.Set(ctx, cooldownKey, "1", emailCodeCooldown).Err()
}
