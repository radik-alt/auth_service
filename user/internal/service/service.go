package service

import (
	"context"
	"errors"

	"github.com/radik-alt/auth_service/user/internal/model"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
	ErrPhoneTaken    = errors.New("phone already taken")

	// ErrInvalidCredentials — и для неверного пароля, и для несуществующего email:
	// по ответу нельзя узнать, зарегистрирован ли email
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserBlocked        = errors.New("user is blocked")
)

type UserService interface {
	CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error)
	UpdateUser(ctx context.Context, userID string, params model.UpdateUserParams) (*model.User, error)
	VerifyCredentials(ctx context.Context, email, password string) (*model.User, error)
}
