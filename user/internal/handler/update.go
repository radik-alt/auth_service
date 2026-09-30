package handler

import (
	"context"
	"strings"

	"github.com/radik-alt/auth_service/user/internal/handler/convertor"
	"github.com/radik-alt/auth_service/user/internal/model"
	userv1 "github.com/radik-alt/auth_service/user/pkg/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *userHandler) UpdateUser(
	ctx context.Context,
	request *userv1.UpdateUserRequest,
) (*userv1.User, error) {
	if request.GetUser() == nil {
		return nil, requiredField("user")
	}
	userID := request.GetUser().GetUserId()
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	caller, err := callerID(ctx)
	if err != nil {
		return nil, err
	}
	if caller != userID {
		return nil, status.Error(codes.PermissionDenied, "cannot modify another user")
	}

	params, err := buildUpdateParams(request.GetUser(), request.GetUpdateMask().GetPaths())
	if err != nil {
		return nil, err
	}

	user, err := handler.service.UpdateUser(ctx, userID, params)
	if err != nil {
		return nil, toStatusError(err)
	}

	return convertor.MapUserToProto(user), nil
}

func buildUpdateParams(user *userv1.User, paths []string) (model.UpdateUserParams, error) {
	var params model.UpdateUserParams

	if len(paths) == 0 {
		return params, requiredField("update_mask")
	}

	for _, path := range paths {
		switch path {
		case "username":
			username := strings.TrimSpace(user.GetUsername())
			if err := validateUsername(username); err != nil {
				return params, err
			}
			params.Username = &username

		case "phone":
			// phone в маске, но не передан — удаляем телефон
			phone := ""
			if user.Phone != nil {
				phone = strings.TrimSpace(user.GetPhone())
				if err := validatePhone(&phone); err != nil {
					return params, err
				}
			}
			params.Phone = &phone

		case "avatar_url":
			avatarURL := strings.TrimSpace(user.GetAvatarUrl())
			if err := validateAvatarURL(avatarURL); err != nil {
				return params, err
			}
			params.AvatarURL = &avatarURL

		case "bio":
			bio := strings.TrimSpace(user.GetBio())
			if err := validateBio(bio); err != nil {
				return params, err
			}
			params.Bio = &bio

		default:
			return params, invalidField("update_mask", "field "+path+" cannot be updated")
		}
	}

	return params, nil
}
