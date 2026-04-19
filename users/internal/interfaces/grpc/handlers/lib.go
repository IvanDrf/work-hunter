package handlers

import (
	"encoding/json"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertCreateProfileResponseToDto(req *user_api.CreateProfileRequest) *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		ID:          req.UserId,
		Username:    req.Username,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}
}

func convertUpdateProfileResponseToDto(req *user_api.UpdateProfileRequest) (*dto.UpdateUserRequest, error) {
	dto := &dto.UpdateUserRequest{
		ID: req.UserId,
	}

	if req.FirstName != nil {
		dto.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		dto.LastName = *req.LastName
	}

	if req.PhoneNumber != nil {
		dto.PhoneNumber = *req.PhoneNumber
	}

	if req.AvatarUrl != nil {
		dto.AvatarURL = *req.AvatarUrl
	}

	data, err := json.Marshal(req.Metadata)
	if err != nil {
		return nil, err
	}

	dto.Metadata = data

	return dto, nil
}

func convertUserResponseToUserProfile(user *dto.UserResponse) *user_api.UserProfile {
	return &user_api.UserProfile{
		Id:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Avatar_URL:  user.AvatarURL,
		Status:      user_api.UserStatus(user_api.UserStatus_value[user.Status]),
		Role:        user_api.UserRole(user_api.UserRole_value[user.Role]),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
		Metadata:    user.Metadata,
	}
}
