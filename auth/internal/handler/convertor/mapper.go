package convertor

import (
	"github.com/radik-alt/auth_service/auth/internal/model"
	pb "github.com/radik-alt/auth_service/auth/pkg/auth/v1"
)

func MapAuthTokensToProto(tokens *model.TokenPair) *pb.AuthResponse {
	if tokens == nil {
		return &pb.AuthResponse{}
	}

	return &pb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    "Bearer",
	}
}
