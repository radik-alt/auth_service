package handler

import (
	"context"

	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) RefreshToken(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (*pb.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, requiredField("refresh_token")
	}

	token, err := handler.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, toStatusError(err)
	}

	return convertor.MapAuthTokensToProto(token), nil
}
