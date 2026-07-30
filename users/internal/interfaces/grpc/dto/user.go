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
	Page     int
	PageSize int
	Filter   *UserFilter
	SortBy   string
	SortDesc bool
}

type UserFilter struct {
	Status      string
	Role        string
	SearchQuery string
}

type ListUsersResponse struct {
	Users      []*UserResponse
	TotalCount int32
	Page       int
	PageSize   int
	TotalPages int32
}

// DTO for updating user status
type UpdateUserStatusRequest struct {
	ID     string
	Status string
}
