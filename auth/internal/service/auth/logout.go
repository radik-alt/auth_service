package auth

import (
	"context"
	"errors"
	"time"

	"github.com/radik-alt/auth_service/auth/internal/repository"
)

func (s *authService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return ErrInvalidToken
	}

	if ttl := time.Until(claims.ExpiresAt.Time); ttl > 0 {
		if err := s.blacklist.Add(ctx, claims.ID, ttl); err != nil {
			return err
		}
	}

	if refreshToken != "" {
		if err := s.refreshRepo.DeleteByTokenHash(ctx, hashToken(refreshToken)); err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}
	}

	return nil
}
