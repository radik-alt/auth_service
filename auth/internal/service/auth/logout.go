package auth

import (
	"context"
	"errors"
	"time"

	"github.com/radik-alt/auth_service/auth/internal/repository"
)

func (s *authService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	claims, err := s.parseAccessTokenIgnoringExpiry(accessToken)
	if err != nil {
		return ErrInvalidToken
	}

	if refreshToken != "" {
		hash := hashToken(refreshToken)

		stored, err := s.refreshRepo.GetByTokenHash(ctx, hash)
		switch {
		case err == nil:
			// refresh чужого пользователя: ничего не удаляем
			if stored.UserID != claims.UserID {
				return ErrInvalidToken
			}
			if err := s.refreshRepo.DeleteByTokenHash(ctx, hash); err != nil {
				return err
			}
		case !errors.Is(err, repository.ErrNotFound):
			return err
		}
	}

	if claims.ExpiresAt == nil {
		return nil
	}
	return s.revoke(ctx, claims.ID, claims.ExpiresAt.Time)
}

// revoke кладёт access-токен в blacklist до его естественного истечения.
func (s *authService) revoke(ctx context.Context, jti string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if jti == "" || ttl <= 0 {
		return nil
	}
	return s.blacklist.Add(ctx, jti, ttl)
}
