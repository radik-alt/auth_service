package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/repository"
)

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	stored, err := s.refreshRepo.ConsumeByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("consume refresh token: %w", err)
	}

	if stored.IsExpired() {
		return nil, ErrInvalidToken
	}

	accessToken, jti, accessExpiresAt, err := s.generateAccessToken(stored.UserID)
	if err != nil {
		s.restore(ctx, stored)
		return nil, err
	}

	newRefreshToken := newRefreshTokenValue()
	now := time.Now()

	if err := s.refreshRepo.SaveByTokenHash(ctx, &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    stored.UserID,
		TokenHash: hashToken(newRefreshToken),
		ExpiresAt: now.Add(s.refreshTTL),
		CreatedAt: now,

		AccessJTI:       jti,
		AccessExpiresAt: accessExpiresAt,
	}); err != nil {
		s.restore(ctx, stored)
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	// старый access больше не нужен; сбой отзыва не должен ломать уже выданную пару
	if err := s.revoke(ctx, stored.AccessJTI, stored.AccessExpiresAt); err != nil {
		log.Printf("revoke previous access token: %v", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    uint64(s.accessTTL.Seconds()),
	}, nil
}

// restore возвращает погашенный refresh-токен, чтобы клиент мог повторить запрос, а не потерял сессию.
func (s *authService) restore(ctx context.Context, stored *model.RefreshToken) {
	if err := s.refreshRepo.SaveByTokenHash(ctx, stored); err != nil {
		log.Printf("restore refresh token: %v", err)
	}
}
