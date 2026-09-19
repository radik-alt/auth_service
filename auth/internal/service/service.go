package service

import (
	"context"
	"errors"

	"github.com/radik-alt/auth_service/auth/internal/model"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type AuthService interface {
	LogIn(ctx context.Context, userID uint64) (*model.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	LogoutAll(ctx context.Context, accessToken string) error
}
