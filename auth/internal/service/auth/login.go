package auth

import (
	"context"

	"github.com/radik-alt/auth_service/auth/internal/model"
)

func (s *authService) LogIn(ctx context.Context, userID uint64) (*model.TokenPair, error) {
	return s.issueTokenPair(ctx, userID)
}
