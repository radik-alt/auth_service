package handler

import (
	"context"

	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	"github.com/radik-alt/auth_service/auth/proto"
)

func (handler *authHandler) LogIn(
	ctx context.Context,
	request *proto.LogInRequest,
) (*proto.AuthResponse, error) {
	tokens, err := handler.service.LogIn(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	return convertor.MapAuthTokensToProto(tokens), nil
}
