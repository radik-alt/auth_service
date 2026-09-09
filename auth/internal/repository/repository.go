package repository

import (
	"context"
	"errors"
	"time"

	"github.com/radik-alt/auth_service/auth/internal/model"
)

var ErrNotFound = errors.New("refresh token not found")

type TokenRepository interface {
	SaveByTokenHash(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}

// TokenBlacklist хранит отозванные access-токены (по их jti) до истечения их TTL.
// Позволяет сервисам, проверяющим JWT самостоятельно, мгновенно узнать об отзыве через Logout.
type TokenBlacklist interface {
	Add(ctx context.Context, jti string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}
