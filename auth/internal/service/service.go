package service

import (
	"context"

	"github.com/radik-alt/auth_service/auth/internal/model"
)

type AuthService interface {
	LogIn(ctx context.Context, userID string) (*model.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
}
