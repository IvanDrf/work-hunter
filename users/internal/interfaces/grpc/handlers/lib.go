package handlers

import (
	"github.com/IvanDrf/work-hunter/pkg/common"
	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertCreateProfileResponseToDto(req *user_api.CreateProfileRequest) *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		ID:          req.UserId,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		CompanyName: req.CompanyName,
		Verificated: req.Verificated,
	}
}

func convertUpdateProfileRequestToDto(req *user_api.UpdateProfileRequest) (*dto.UpdateUserRequest, error) {
	dto := &dto.UpdateUserRequest{
		ID: req.UserId,
	}

	if req.FirstName != nil {
		dto.FirstName = req.FirstName
	}

	if req.LastName != nil {
		dto.LastName = req.LastName
	}

	if req.CompanyName != nil {
		dto.CompanyName = req.CompanyName
	}

	if req.Verificated != nil {
		dto.Verificated = req.Verificated
	}
	return dto, nil
}

func convertUserResponseToUserProfile(user *dto.UserResponse) *user_api.UserProfile {
	return &user_api.UserProfile{
		Id:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		CompanyName: user.CompanyName,
		Status:      user_api.UserStatus(user_api.UserStatus_value[user.Status]),
		Role:        common.UserRole(common.UserRole_value[user.Role]),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
		Verificated: user.Verificated,
	}
}

func convertListDtoToListResp(dto *dto.ListUsersResponse) *user_api.ListUsersResponse {
	resp := &user_api.ListUsersResponse{
		Users:      make([]*user_api.UserProfile, len(dto.Users)),
		TotalCount: dto.TotalCount,
	}

	for _, val := range dto.Users {
		resp.Users = append(resp.Users, convertUserResponseToUserProfile(val))
	}

	return resp
}
