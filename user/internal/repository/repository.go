package repository

import (
	"context"
	"errors"

	"github.com/radik-alt/auth_service/user/internal/model"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
	ErrPhoneTaken    = errors.New("phone already taken")
)

type UserRepository interface {
	CreateUser(ctx context.Context, param model.CreateUserParams) (*model.User, error)
	UpdateUser(ctx context.Context, userID string, params model.UpdateUserParams) (*model.User, error)
	// GetByEmail ищет среди неудалённых; email должен быть уже в нижнем регистре. ErrNotFound — не найден.
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
