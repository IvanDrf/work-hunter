package models

import (
	"time"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/rules"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/interfaces/grpc/dto"
)

type User struct {
	ID          string `db:"id" json:"id"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name" json:"last_name"`
	PhoneNumber string `db:"phone_number" json:"phone_number"`
	AvatarURL   string `db:"avatar_url" json:"avatar_url"`

	Status rules.UserStatus `db:"status"`
	Role   rules.UserRole   `db:"role"`

	Metadata map[string]string `db:"metadata"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updatet_at"`
}

func NewUser(req *dto.CreateUserRequest) *User {
	now := time.Now()
	return &User{
		ID:          req.ID,
		Username:    req.Username,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,

		Status: rules.UserStatusActive,
		Role:   rules.UserRoleUser,

		Metadata: make(map[string]string),

		CreatedAt: now,
		UpdatedAt: now,
	}
}
