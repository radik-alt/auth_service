package handler

import (
	"context"

	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) LogIn(
	ctx context.Context,
	request *pb.LogInRequest,
) (*pb.AuthResponse, error) {
	tokens, err := handler.service.LogIn(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	return convertor.MapAuthTokensToProto(tokens), nil
}
