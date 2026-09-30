package auth

import (
	"context"
	"fmt"

	"github.com/radik-alt/auth_service/auth/internal/service"
)

func (s *authService) ValidateToken(ctx context.Context, accessToken string) (string, error) {
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return "", service.ErrInvalidToken
	}

	revoked, err := s.blacklist.IsBlacklisted(ctx, claims.ID)
	if err != nil {
		return "", fmt.Errorf("check blacklist: %w", err)
	}
	if revoked {
		return "", service.ErrInvalidToken
	}

	return claims.UserID, nil
}
