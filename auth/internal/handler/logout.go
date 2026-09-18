package handler

import (
	"context"

	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) Logout(
	ctx context.Context,
	req *pb.LogoutRequest,
) (*pb.LogoutResponse, error) {

	err := handler.service.Logout(ctx, req.AccessToken, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &pb.LogoutResponse{}, nil
}
