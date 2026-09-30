package handler

import (
	"errors"
	"log"

	"github.com/radik-alt/auth_service/user/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrEmailTaken),
		errors.Is(err, service.ErrUsernameTaken),
		errors.Is(err, service.ErrPhoneTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrUserBlocked):
		return status.Error(codes.PermissionDenied, err.Error())
	}

	log.Printf("internal error: %v", err)
	return status.Error(codes.Internal, "internal error")
}

func requiredField(name string) error {
	return status.Errorf(codes.InvalidArgument, "%s is required", name)
}

func invalidField(name, reason string) error {
	return status.Errorf(codes.InvalidArgument, "%s: %s", name, reason)
}
