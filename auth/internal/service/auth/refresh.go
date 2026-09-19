package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/repository"
	"github.com/radik-alt/auth_service/auth/internal/service"
)

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	stored, err := s.refreshRepo.ConsumeByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrInvalidToken
		}
		return nil, fmt.Errorf("consume refresh token: %w", err)
	}

	if stored.IsExpired() {
		return nil, service.ErrInvalidToken
	}

	pair, err := s.issueTokenPair(ctx, stored.UserID)
	if err != nil {
		s.restore(ctx, stored)
		return nil, err
	}

	// старый access больше не нужен; сбой отзыва не должен ломать уже выданную пару
	if err := s.revoke(ctx, stored.AccessJTI, stored.AccessExpiresAt); err != nil {
		log.Printf("revoke previous access token: %v", err)
	}

	return pair, nil
}

// restore возвращает погашенный refresh-токен, чтобы клиент мог повторить запрос, а не потерял сессию.
func (s *authService) restore(ctx context.Context, stored *model.RefreshToken) {
	if err := s.refreshRepo.SaveByTokenHash(ctx, stored); err != nil {
		log.Printf("restore refresh token: %v", err)
	}
}
