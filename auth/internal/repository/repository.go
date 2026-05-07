package repository

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/model"
)

type TokenRepository interface {
	SaveByTokenHash(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, userID string) (string, error)
	DeleteByTokenHash(ctx context.Context, token string) error
}
