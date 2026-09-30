package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/radik-alt/auth_service/user/internal/model"
	"github.com/radik-alt/auth_service/user/internal/repository"
	"github.com/radik-alt/auth_service/user/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.CreateUser(ctx, model.CreateUserParams{
		Email:        strings.ToLower(input.Email),
		Username:     input.Username,
		Phone:        input.Phone,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, mapRepoError(err)
	}

	return user, nil
}

func mapRepoError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return service.ErrUserNotFound
	case errors.Is(err, repository.ErrEmailTaken):
		return service.ErrEmailTaken
	case errors.Is(err, repository.ErrUsernameTaken):
		return service.ErrUsernameTaken
	case errors.Is(err, repository.ErrPhoneTaken):
		return service.ErrPhoneTaken
	default:
		return err
	}
}
