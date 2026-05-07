package handler

import "github.com/radik-alt/auth_service/user/proto"

type userHandler struct {
	*proto.UnimplementedUserServiceServer
}

func NewUserHandler() userHandler {
	return userHandler{}
}
