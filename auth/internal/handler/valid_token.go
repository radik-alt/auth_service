package handler

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	"github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *authHandler) ValidateToken(
	ctx context.Context,
	proto *proto.ValidateTokenRequest,
) (*proto.ValidateTokenResponse, error) {

	isValid, err := handler.service.ValidateToken(ctx, proto.Token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "method invalidate")
	}

	return convertor.MapValidateToken(isValid), nil
}
