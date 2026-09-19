package handler

import (
	"context"

	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) Logout(
	ctx context.Context,
	req *pb.LogoutRequest,
) (*pb.LogoutResponse, error) {
	if req.AccessToken == "" {
		return nil, requiredField("access_token")
	}

	if err := handler.service.Logout(ctx, req.AccessToken, req.RefreshToken); err != nil {
		return nil, toStatusError(err)
	}

	return &pb.LogoutResponse{}, nil
}
