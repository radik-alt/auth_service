package auth

import "context"

// LogoutAll завершает все сессии пользователя access-токена (выход со всех устройств).
func (s *authService) LogoutAll(ctx context.Context, accessToken string) error {
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return ErrInvalidToken
	}

	sessions, err := s.refreshRepo.DeleteAllByUser(ctx, claims.UserID)
	if err != nil {
		return err
	}

	if claims.ExpiresAt != nil {
		if err := s.revoke(ctx, claims.ID, claims.ExpiresAt.Time); err != nil {
			return err
		}
	}
	for _, session := range sessions {
		if err := s.revoke(ctx, session.AccessJTI, session.AccessExpiresAt); err != nil {
			return err
		}
	}

	return nil
}
