package handler

import (
	"context"

	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func (handler *authHandler) ValidateToken(
	ctx context.Context,
	req *pb.ValidateTokenRequest,
) (*pb.ValidateTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, requiredField("access_token")
	}

	userID, err := handler.service.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &pb.ValidateTokenResponse{UserId: userID}, nil
}
