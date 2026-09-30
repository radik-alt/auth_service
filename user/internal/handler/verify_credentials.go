package handler

import (
	"context"
	"strings"

	"github.com/radik-alt/auth_service/user/internal/handler/convertor"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
)

func (handler *userHandler) VerifyCredentials(
	ctx context.Context,
	request *userv1.VerifyCredentialsRequest,
) (*userv1.VerifyCredentialsResponse, error) {
	email := strings.TrimSpace(request.GetEmail())
	if email == "" {
		return nil, requiredField("email")
	}
	if request.GetPassword() == "" {
		return nil, requiredField("password")
	}
	if len(request.GetPassword()) > maxPasswordLen {
		return nil, invalidField("password", "must be at most 72 bytes")
	}

	user, err := handler.service.VerifyCredentials(ctx, email, request.GetPassword())
	if err != nil {
		return nil, toStatusError(err)
	}

	return convertor.MapVerifyCredentialsToProto(user), nil
}
