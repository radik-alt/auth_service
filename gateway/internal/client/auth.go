package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	authpb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrUnavailable  = errors.New("auth service unavailable")
)

// requestTimeout не даёт запросу висеть, если auth тормозит
const requestTimeout = 3 * time.Second

type AuthClient struct {
	api authpb.AuthServiceClient
}

func NewAuthClient(conn grpc.ClientConnInterface) *AuthClient {
	return &AuthClient{api: authpb.NewAuthServiceClient(conn)}
}

// ValidateToken возвращает user_id владельца access-токена.
func (c *AuthClient) ValidateToken(ctx context.Context, accessToken string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := c.api.ValidateToken(ctx, &authpb.ValidateTokenRequest{AccessToken: accessToken})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated, codes.InvalidArgument:
			return "", ErrInvalidToken
		case codes.Unavailable, codes.DeadlineExceeded:
			return "", fmt.Errorf("validate token: %w: %s", ErrUnavailable, status.Convert(err).Message())
		default:
			return "", fmt.Errorf("validate token: %w", err)
		}
	}

	return resp.GetUserId(), nil
}
