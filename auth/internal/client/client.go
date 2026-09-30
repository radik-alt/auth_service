package client

import (
	"context"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserBlocked        = errors.New("user is blocked")
	ErrUnavailable        = errors.New("unavailable")
)

// UserClient — обращения auth к user-сервису.
type UserClient interface {
	// VerifyCredentials возвращает user_id (UUID) при верных email и пароле.
	// ErrInvalidCredentials — неверные данные, ErrUserBlocked — пользователь заблокирован,
	// ErrUnavailable — user-сервис не ответил.
	VerifyCredentials(ctx context.Context, email, password string) (string, error)
}
