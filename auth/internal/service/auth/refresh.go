package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/repository"
)

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	tokenHash := hashToken(refreshToken)

	stored, err := s.refreshRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if stored.IsExpired() {
		_ = s.refreshRepo.DeleteByTokenHash(ctx, tokenHash)
		return nil, ErrInvalidToken
	}

	if err := s.refreshRepo.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("delete old refresh token: %w", err)
	}

	accessToken, _, _, err := s.generateAccessToken(stored.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken := newRefreshTokenValue()
	now := time.Now()
	expiresAt := now.Add(s.refreshTTL)

	if err := s.refreshRepo.SaveByTokenHash(ctx, &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    stored.UserID,
		TokenHash: hashToken(newRefreshToken),
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, nil
}
