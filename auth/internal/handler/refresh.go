package handler

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	"github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *authHandler) RefreshToken(
	ctx context.Context,
	proto *proto.RefreshTokenRequest,
) (*proto.AuthResponse, error) {

	token, err := handler.service.RefreshToken(ctx, proto.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, "method RefreshToken not implemented")
	}

	return convertor.MapAuthTokensToProto(token, true, ""), nil
}
