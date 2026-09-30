package handler

import (
	"context"
	"strings"

	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) LogIn(
	ctx context.Context,
	request *pb.LogInRequest,
) (*pb.AuthResponse, error) {
	email := strings.TrimSpace(request.GetEmail())
	if email == "" {
		return nil, requiredField("email")
	}
	if request.GetPassword() == "" {
		return nil, requiredField("password")
	}

	tokens, err := handler.service.LogIn(ctx, email, request.GetPassword())
	if err != nil {
		return nil, toStatusError(err)
	}

	return convertor.MapAuthTokensToProto(tokens), nil
}
