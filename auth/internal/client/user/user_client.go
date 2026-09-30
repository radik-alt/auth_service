package user

import (
	"context"
	"fmt"
	"time"

	"github.com/radik-alt/auth_service/auth/internal/client"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// requestTimeout не даёт логину висеть, если user-сервис тормозит
const requestTimeout = 3 * time.Second

type userClient struct {
	api userv1.UserServiceClient
}

func NewUserClient(conn grpc.ClientConnInterface) client.UserClient {
	return &userClient{api: userv1.NewUserServiceClient(conn)}
}

func (c *userClient) VerifyCredentials(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := c.api.VerifyCredentials(ctx, &userv1.VerifyCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", mapStatusError(err)
	}

	return resp.GetUserId(), nil
}

func mapStatusError(err error) error {
	switch status.Code(err) {
	// InvalidArgument — например, пароль длиннее допустимого: для клиента это те же неверные данные
	case codes.Unauthenticated, codes.InvalidArgument:
		return client.ErrInvalidCredentials
	case codes.PermissionDenied:
		return client.ErrUserBlocked
	case codes.Unavailable, codes.DeadlineExceeded:
		return fmt.Errorf("verify credentials: %w: %s", client.ErrUnavailable, status.Convert(err).Message())
	default:
		return fmt.Errorf("verify credentials: %w", err)
	}
}
