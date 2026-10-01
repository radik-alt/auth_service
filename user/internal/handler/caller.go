package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const callerIDMetadataKey = "x-user-id"

// callerID возвращает id того, кто делает запрос.
func callerID(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md.Get(callerIDMetadataKey)
	if len(values) == 0 || values[0] == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return values[0], nil
}
