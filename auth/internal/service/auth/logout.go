package auth

import (
	"context"
	"time"
)

// TODO: LogoutRequest несёт только access-токен, поэтому refresh-токен той же
// сессии здесь не отзывается и остаётся рабочим до истечения своего TTL.
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return err
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return s.blacklist.Add(ctx, claims.ID, ttl)
}
