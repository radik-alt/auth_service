package auth

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/model"
)

func (a authService) ValidateToken(ctx context.Context, token string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a authService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	//TODO implement me
	panic("implement me")
}

func (a authService) Logout(ctx context.Context, refreshToken string) error {
	//TODO implement me
	panic("implement me")
}
