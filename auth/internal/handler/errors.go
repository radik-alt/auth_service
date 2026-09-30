package handler

import (
	"errors"
	"log"

	"github.com/radik-alt/auth_service/auth/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, "invalid or expired token")
	case errors.Is(err, service.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid email or password")
	case errors.Is(err, service.ErrUserBlocked):
		return status.Error(codes.PermissionDenied, "user is blocked")
	case errors.Is(err, service.ErrUserServiceUnavailable):
		log.Printf("login: %v", err)
		return status.Error(codes.Unavailable, "service temporarily unavailable")
	}

	log.Printf("internal error: %v", err)
	return status.Error(codes.Internal, "internal error")
}

func requiredField(name string) error {
	return status.Errorf(codes.InvalidArgument, "%s is required", name)
}
