package repository

import (
	"context"
	"errors"

	"github.com/radik-alt/auth_service/auth/internal/model"
)

var ErrNotFound = errors.New("refresh token not found")

type TokenRepository interface {
	SaveByTokenHash(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}
