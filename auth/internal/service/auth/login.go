package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/radik-alt/auth_service/auth/internal/client"
	"github.com/radik-alt/auth_service/auth/internal/model"
	"github.com/radik-alt/auth_service/auth/internal/service"
)

// LogIn проверяет пароль в user-сервисе и выдаёт токены на его user_id.
func (s *authService) LogIn(ctx context.Context, email, password string) (*model.TokenPair, error) {
	userID, err := s.userClient.VerifyCredentials(ctx, email, password)
	if err != nil {
		return nil, mapClientError(err)
	}

	return s.issueTokenPair(ctx, userID)
}

func mapClientError(err error) error {
	switch {
	case errors.Is(err, client.ErrInvalidCredentials):
		return service.ErrInvalidCredentials
	case errors.Is(err, client.ErrUserBlocked):
		return service.ErrUserBlocked
	case errors.Is(err, client.ErrUnavailable):
		// исходную ошибку сохраняем для лога
		return fmt.Errorf("%w: %v", service.ErrUserServiceUnavailable, err)
	default:
		return err
	}
}
