package handler

import (
	"context"
	"github.com/radik-alt/auth_service/auth/internal/handler/convertor"
	"github.com/radik-alt/auth_service/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *authHandler) Logout(
	ctx context.Context,
	req *proto.LogoutRequest,
) (*proto.LogoutResponse, error) {

	err := handler.service.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return convertor.MapLogoutResponse(false, "error logout"), status.Error(codes.Unimplemented, "method RefreshToken not implemented")
	}

	return convertor.MapLogoutResponse(true, "success"), nil
}
