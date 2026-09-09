package handler

import (
	"context"

	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	"github.com/radik-alt/auth_service/auth/proto"
)

func (handler *authHandler) RefreshToken(
	ctx context.Context,
	proto *proto.RefreshTokenRequest,
) (*proto.AuthResponse, error) {

	token, err := handler.service.RefreshToken(ctx, proto.RefreshToken)
	if err != nil {
		return nil, err
	}

	return convertor.MapAuthTokensToProto(token), nil
}
