package user

import (
	"context"

	"github.com/radik-alt/auth_service/user/internal/model"
)

func (s *userService) UpdateUser(ctx context.Context, userID string, params model.UpdateUserParams) (*model.User, error) {
	user, err := s.userRepo.UpdateUser(ctx, userID, params)
	if err != nil {
		return nil, mapRepoError(err)
	}

	return user, nil
}
