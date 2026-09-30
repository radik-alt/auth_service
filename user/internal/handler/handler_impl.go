package handler

import (
	"github.com/radik-alt/auth_service/user/internal/service"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
)

type userHandler struct {
	userv1.UnimplementedUserServiceServer
	service service.UserService
}

func NewUserHandler(s service.UserService) *userHandler {
	return &userHandler{
		service: s,
	}
}
