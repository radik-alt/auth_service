package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/radik-alt/auth_service/auth/internal/repository"
	"github.com/redis/go-redis/v9"
)

const blacklistKeyPrefix = "access_blacklist:"

type blacklistRepository struct {
	client *redis.Client
}

func NewBlacklistRepo(client *redis.Client) repository.TokenBlacklist {
	return &blacklistRepository{client: client}
}

func (r *blacklistRepository) Add(ctx context.Context, jti string, ttl time.Duration) error {
	if err := r.client.Set(ctx, blacklistKeyPrefix+jti, "1", ttl).Err(); err != nil {
		return fmt.Errorf("blacklist access token: %w", err)
	}
	return nil
}

func (r *blacklistRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	n, err := r.client.Exists(ctx, blacklistKeyPrefix+jti).Result()
	if err != nil {
		return false, fmt.Errorf("check blacklist: %w", err)
	}
	return n > 0, nil
}
