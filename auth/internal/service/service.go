package service

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/model"
)

type AuthService interface {
	GenerateTokens(ctx context.Context, userID, email string) (*model.TokenPair, error)
	ValidateToken(ctx context.Context, accessToken string) (bool, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
}
