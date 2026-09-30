package service

import (
	"context"
	"errors"

	"github.com/radik-alt/auth_service/auth/internal/model"
)

var (
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrUserBlocked            = errors.New("user is blocked")
	ErrUserServiceUnavailable = errors.New("user service unavailable")
)

type AuthService interface {
	LogIn(ctx context.Context, email, password string) (*model.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	LogoutAll(ctx context.Context, accessToken string) error
}
