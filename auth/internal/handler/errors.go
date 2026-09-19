package handler

import (
	"errors"
	"log"

	"github.com/radik-alt/auth_service/auth/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatusError(err error) error {
	if errors.Is(err, service.ErrInvalidToken) {
		return status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	log.Printf("internal error: %v", err)
	return status.Error(codes.Internal, "internal error")
}

func requiredField(name string) error {
	return status.Errorf(codes.InvalidArgument, "%s is required", name)
}
