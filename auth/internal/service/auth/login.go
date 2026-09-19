package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/radik-alt/auth_service/auth/internal/model"
)

func (s *authService) LogIn(ctx context.Context, userID uint64) (*model.TokenPair, error) {
	accessToken, jti, accessExpiresAt, err := s.generateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	isBlack, err := s.blacklist.IsBlacklisted(ctx, jti)
	if err != nil {
		return nil, err
	}

	if isBlack {
		return nil, ErrTokenExpired
	}

	refreshToken := newRefreshTokenValue()
	now := time.Now()
	expiresAt := now.Add(s.refreshTTL)

	err = s.refreshRepo.SaveByTokenHash(ctx, &model.RefreshToken{
		ID:              uuid.NewString(),
		UserID:          userID,
		TokenHash:       hashToken(refreshToken),
		ExpiresAt:       expiresAt,
		CreatedAt:       now,
		AccessJTI:       jti,
		AccessExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uint64(s.accessTTL.Seconds()),
	}, nil
}
