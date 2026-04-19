package dto

import (
	"time"
)

type UserResponse struct {
	ID          string
	Username    string
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string

	Status string
	Role   string

	CreatedAt time.Time
	UpdatedAt time.Time

	Metadata map[string]string
}

// DTO for creating user
type CreateUserRequest struct {
	ID          string
	Username    string
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
}

// DTO for updating user
type UpdateUserRequest struct {
	ID          string
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string
	Metadata    []byte
}

// DTO for listing users
type ListUsersRequest struct {
	PageSize    int32
	Status      string
	Role        string
	SearchQuery string
	SortBy      string
	Offset      int32
}

type ListUsersResponse struct {
	Users      []*UserResponse
	TotalCount int32
}

// DTO for updating user status
type UpdateUserStatusRequest struct {
	ID     string
	Status string
}
