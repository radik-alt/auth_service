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

const (
	refreshTokenKeyPrefix = "refresh_token:"
	userSessionsKeyPrefix = "user_sessions:"
)

// userSessionsKey — множество хэшей refresh-токенов пользователя (индекс «пользователь -> сессии»).
func userSessionsKey(userID string) string {
	return userSessionsKeyPrefix + userID
}

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

	index := userSessionsKey(token.UserID)

	pipe := r.client.TxPipeline()
	pipe.Set(ctx, refreshTokenKeyPrefix+token.TokenHash, data, ttl)
	pipe.SAdd(ctx, index, token.TokenHash)

	pipe.ExpireNX(ctx, index, ttl)
	pipe.ExpireGT(ctx, index, ttl)
	if _, err := pipe.Exec(ctx); err != nil {
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
	return decodeToken(data)
}

func (r *refreshTokenRepository) ConsumeByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	data, err := r.client.GetDel(ctx, refreshTokenKeyPrefix+tokenHash).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("consume refresh token: %w", err)
	}

	token, err := decodeToken(data)
	if err != nil {
		return nil, err
	}

	_ = r.client.SRem(ctx, userSessionsKey(token.UserID), tokenHash).Err()
	return token, nil
}

func (r *refreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	if _, err := r.ConsumeByTokenHash(ctx, tokenHash); err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteAllByUser(ctx context.Context, userID string) ([]*model.RefreshToken, error) {
	index := userSessionsKey(userID)

	// Снимок индекса и его удаление в одной транзакции: сессия, созданная позже, попадёт уже в новый индекс.
	pipe := r.client.TxPipeline()
	members := pipe.SMembers(ctx, index)
	pipe.Del(ctx, index)
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("read user sessions: %w", err)
	}

	hashes := members.Val()
	if len(hashes) == 0 {
		return nil, nil
	}

	getPipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(hashes))
	for i, hash := range hashes {
		cmds[i] = getPipe.GetDel(ctx, refreshTokenKeyPrefix+hash)
	}
	// redis.Nil от Exec означает «часть ключей уже истекла», это нормально
	if _, err := getPipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("delete user sessions: %w", err)
	}

	tokens := make([]*model.RefreshToken, 0, len(cmds))
	for _, cmd := range cmds {
		data, err := cmd.Bytes()
		if errors.Is(err, redis.Nil) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("delete user sessions: %w", err)
		}
		token, err := decodeToken(data)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func decodeToken(data []byte) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unmarshal refresh token: %w", err)
	}
	return &token, nil
}
