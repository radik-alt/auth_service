package handler

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	"github.com/radik-alt/auth_service/auth/proto"
)

func (handler *authHandler) GenerateTokens(
	ctx context.Context,
	request *proto.GenerateTokensRequest,
) (*proto.AuthResponse, error) {
	tokens, err := handler.service.GenerateTokens(ctx, request.UserId, request.Email)
	if err != nil {
		return convertor.MapAuthTokensToProto(nil, false, "error generate tokens"), err
	}

	return convertor.MapAuthTokensToProto(tokens, true, ""), nil
}
