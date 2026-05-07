package convertor

import (
	"github.com/radik-alt/auth_service/auth/internal/model"
	pb "github.com/radik-alt/auth_service/auth/proto"
)

func MapAuthTokensToProto(tokens *model.TokenPair, success bool, message string) *pb.AuthResponse {
	if tokens == nil {
		return &pb.AuthResponse{
			Success: success,
			Message: message,
		}
	}

	return &pb.AuthResponse{
		Success:      success,
		Message:      message,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    "Bearer",
	}
}

func MapValidateToken(isValid bool) *pb.ValidateTokenResponse {
	return &pb.ValidateTokenResponse{
		Valid: isValid,
	}
}

func MapLogoutResponse(
	isSuccess bool,
	message string,
) *pb.LogoutResponse {
	return &pb.LogoutResponse{
		Success: isSuccess,
		Message: message,
	}
}
