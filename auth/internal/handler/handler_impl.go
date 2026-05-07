package handler

import (
	"github.com/radik-alt/auth_service/auth/internal/service"
	pb "github.com/radik-alt/auth_service/auth/proto"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{
		service: s,
	}
}
