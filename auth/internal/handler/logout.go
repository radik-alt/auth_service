package handler

import (
	"context"

	"github.com/radik-alt/auth_service/auth/proto"
)

func (handler *authHandler) Logout(
	ctx context.Context,
	req *proto.LogoutRequest,
) (*proto.LogoutResponse, error) {

	err := handler.service.Logout(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &proto.LogoutResponse{}, nil
}
