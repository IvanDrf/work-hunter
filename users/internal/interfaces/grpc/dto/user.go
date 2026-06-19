package dto

import (
	"time"
)

type UserResponse struct {
	ID          string
	Email       string
	FirstName   string
	LastName    string
	CompanyName string

	Status string
	Role   string

	CreatedAt time.Time
	UpdatedAt time.Time

	Verificated bool
}

// DTO for creating user
type CreateUserRequest struct {
	ID          string
	Email       string
	FirstName   string
	LastName    string
	CompanyName string
	Verificated bool
}

// DTO for updating user
type UpdateUserRequest struct {
	ID          string
	FirstName   string
	LastName    string
	CompanyName string
	Verificated bool
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
