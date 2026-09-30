package convertor

import (
	"github.com/radik-alt/auth_service/user/internal/model"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MapUserToProto не переносит PasswordHash и DeletedAt: они не выходят за пределы сервиса.
func MapUserToProto(user *model.User) *userv1.User {
	if user == nil {
		return &userv1.User{}
	}

	return &userv1.User{
		UserId:        user.ID,
		Email:         user.Email,
		Username:      user.Username,
		Phone:         user.Phone,
		AvatarUrl:     user.AvatarURL,
		Bio:           user.Bio,
		Status:        mapStatusToProto(user.Status),
		EmailVerified: user.EmailVerified,
		CreatedAt:     timestamppb.New(user.CreatedAt),
		UpdatedAt:     timestamppb.New(user.UpdatedAt),
	}
}

func MapVerifyCredentialsToProto(user *model.User) *userv1.VerifyCredentialsResponse {
	if user == nil {
		return &userv1.VerifyCredentialsResponse{}
	}

	return &userv1.VerifyCredentialsResponse{
		UserId: user.ID,
		Status: mapStatusToProto(user.Status),
	}
}

func mapStatusToProto(status model.UserStatus) userv1.UserStatus {
	switch status {
	case model.UserStatusActive:
		return userv1.UserStatus_USER_STATUS_ACTIVE
	case model.UserStatusBlocked:
		return userv1.UserStatus_USER_STATUS_BLOCKED
	case model.UserStatusDeleted:
		return userv1.UserStatus_USER_STATUS_DELETED
	default:
		return userv1.UserStatus_USER_STATUS_UNSPECIFIED
	}
}
