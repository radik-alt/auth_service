package handler

import (
	"context"
	"strings"

	"github.com/radik-alt/auth_service/user/internal/handler/convertor"
	"github.com/radik-alt/auth_service/user/internal/model"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
)

func (handler *userHandler) CreateUser(
	ctx context.Context,
	request *userv1.CreateUserRequest,
) (*userv1.User, error) {
	email := strings.TrimSpace(request.GetEmail())
	username := strings.TrimSpace(request.GetUsername())

	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validatePassword(request.GetPassword()); err != nil {
		return nil, err
	}
	if err := validatePhone(request.Phone); err != nil {
		return nil, err
	}

	user, err := handler.service.CreateUser(ctx, model.CreateUserInput{
		Email:    email,
		Username: username,
		Phone:    request.Phone,
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, toStatusError(err)
	}

	return convertor.MapUserToProto(user), nil
}
