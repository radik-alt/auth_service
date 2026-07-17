package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	tokenHash := hashToken(refreshToken)

	stored, err := s.repo.GetByTokenHash(ctx, tokenHash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	// Токен одноразовый: удаляем его при любой попытке использования,
	// даже если он оказался просроченным, чтобы его нельзя было переиспользовать.
	if err := s.repo.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("delete refresh token: %w", err)
	}

	if stored.IsExpired() {
		return nil, status.Error(codes.Unauthenticated, "refresh token expired")
	}

	accessToken, err := s.issueAccessToken(stored.UserID, "")
	if err != nil {
		return nil, err
	}

	rawRefresh, err := s.issueRefreshToken(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, nil
}
