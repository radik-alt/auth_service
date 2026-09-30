package handler

import userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"

type userHandler struct {
	*userv1.UnimplementedUserServiceServer
}

func NewUserHandler() userHandler {
	return userHandler{}
}
