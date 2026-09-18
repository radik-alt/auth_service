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

	token, err := handler.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return convertor.MapAuthTokensToProto(token), nil
}
