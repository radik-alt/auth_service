package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/radik-alt/auth_service/auth/internal/model"
)

/* issueTokenPair выпускает пару access/refresh и сохраняет сессию (общий код LogIn и RefreshToken).
 */
func (s *authService) issueTokenPair(ctx context.Context, userID uint64) (*model.TokenPair, error) {
	accessToken, jti, accessExpiresAt, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken := newRefreshTokenValue()
	now := time.Now()

	if err := s.refreshRepo.SaveByTokenHash(ctx, &model.RefreshToken{
		ID:              uuid.NewString(),
		UserID:          userID,
		TokenHash:       hashToken(refreshToken),
		ExpiresAt:       now.Add(s.refreshTTL),
		CreatedAt:       now,
		AccessJTI:       jti,
		AccessExpiresAt: accessExpiresAt,
	}); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uint64(s.accessTTL.Seconds()),
	}, nil
}
