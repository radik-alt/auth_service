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

func (s *userService) VerifyCredentials(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, service.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("compare password: %w", err)
	}

	if !user.IsActive() {
		return nil, service.ErrUserBlocked
	}

	return user, nil
}
