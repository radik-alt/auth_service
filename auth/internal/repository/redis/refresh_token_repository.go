package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/repository"
	"github.com/redis/go-redis/v9"
)

const refreshTokenKeyPrefix = "refresh_token:"

type refreshTokenRepository struct {
	client *redis.Client
}

func NewRefreshTokenRepo(client *redis.Client) repository.TokenRepository {
	return &refreshTokenRepository{client: client}
}

func (r *refreshTokenRepository) SaveByTokenHash(ctx context.Context, token *model.RefreshToken) error {
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("save refresh token: already expired")
	}

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshal refresh token: %w", err)
	}

	if err := r.client.Set(ctx, refreshTokenKeyPrefix+token.TokenHash, data, ttl).Err(); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	data, err := r.client.Get(ctx, refreshTokenKeyPrefix+tokenHash).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	var token model.RefreshToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unmarshal refresh token: %w", err)
	}
	return &token, nil
}

func (r *refreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	if err := r.client.Del(ctx, refreshTokenKeyPrefix+tokenHash).Err(); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}
