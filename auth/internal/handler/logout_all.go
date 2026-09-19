package handler

import (
	"context"

	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) LogoutAll(
	ctx context.Context,
	req *pb.LogoutAllRequest,
) (*pb.LogoutAllResponse, error) {

	if err := handler.service.LogoutAll(ctx, req.AccessToken); err != nil {
		return nil, err
	}

	return &pb.LogoutAllResponse{}, nil
}
