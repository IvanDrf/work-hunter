package handlers

import (
	"strings"

	"github.com/IvanDrf/work-hunter/pkg/common"
	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	userStatus_DB = map[string]string{"USER_STATUS_ACTIVE": "active", "USER_STATUS_INACTIVE": "inactive", "USER_STATUS_BLOCKED": "blocked", "USER_STATUS_DELETED": "deleted"}
	userRole_PB   = map[string]string{"active": "USER_STATUS_ACTIVE", "inactive": "USER_STATUS_INACTIVE", "blocked": "USER_STATUS_BLOCKED", "deleted": "USER_STATUS_DELETED"}
)

func convertCreateProfileResponseToDto(req *user_api.CreateProfileRequest) *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		ID:          req.UserId,
		Email:       req.Email,
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		CompanyName: req.GetCompanyName(),
		Verificated: req.Verificated,
	}
}

func convertUpdateProfileRequestToDto(req *user_api.UpdateProfileRequest) (*dto.UpdateUserRequest, error) {
	dto := &dto.UpdateUserRequest{
		ID: req.UserId,
	}

	if req.FirstName != nil {
		dto.FirstName = req.GetFirstName()
	}

	if req.LastName != nil {
		dto.LastName = req.GetLastName()
	}

	if req.CompanyName != nil {
		dto.CompanyName = req.GetCompanyName()
	}

	if req.Verificated != nil {
		dto.Verificated = req.GetVerificated()
	}
	return dto, nil
}

func convertUserResponseToUserProfile(user *dto.UserResponse) *user_api.UserProfile {
	u := user_api.UserProfile{
		Id:          user.ID,
		Email:       user.Email,
		Status:      user_api.UserStatus(user_api.UserStatus_value[userRole_PB[user.Status]]),
		Role:        common.UserRole(common.UserRole_value[strings.ToUpper(user.Role)]),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
		Verificated: user.Verificated,
	}

	if user.CompanyName != "" {
		u.CompanyName = &user.CompanyName
	}
	if user.FirstName != "" {
		u.FirstName = &user.FirstName
	}
	if user.LastName != "" {
		u.LastName = &user.LastName
	}

	return &u
}

func convertListDtoToListResp(dto *dto.ListUsersResponse) *user_api.ListUsersResponse {
	resp := &user_api.ListUsersResponse{
		Users:      make([]*user_api.UserProfile, 0, len(dto.Users)),
		TotalCount: dto.TotalCount,
	}

	for _, val := range dto.Users {
		resp.Users = append(resp.Users, convertUserResponseToUserProfile(val))
	}

	return resp
}
