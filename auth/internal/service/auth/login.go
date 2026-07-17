package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	err := s.repo.DeleteByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		return status.Error(codes.NotFound, "refresh token not found")
	}

	return nil
}
