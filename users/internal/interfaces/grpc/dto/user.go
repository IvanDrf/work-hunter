package dto

import (
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/google/uuid"
)

// DTO for creating user
type CreateUserRequest struct {
	ID          uuid.UUID
	Username    string
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
}

// DTO for updating user
type UpdateUserRequest struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string
	Metadata    []byte
}

// DTO for listing users
type ListUsersRequest struct {
	PageSize    int32
	Status      rules.UserStatus
	Role        rules.UserRole
	SearchQuery string
	SortBy      string
	Offset      int32
}

type ListUsersResponse struct {
	Users      []*models.User
	TotalCount int32
	HasNext    bool
}

// DTO for updating user status
type UpdateUserStatusRequest struct {
	ID     string
	Status rules.UserStatus
}
