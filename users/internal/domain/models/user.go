package models

import (
	"encoding/json"
	"time"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/rules"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name" json:"last_name"`
	PhoneNumber string    `db:"phone_number" json:"phone_number"`
	AvatarURL   string    `db:"avatar_url" json:"avatar_url"`

	Status rules.UserStatus `db:"status"`
	Role   rules.UserRole   `db:"role"`

	Metadata json.RawMessage `db:"metadata"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(id uuid.UUID, username, email, firstName, lastName, phoneNumber string) *User {
	now := time.Now()
	return &User{
		ID:          id,
		Username:    username,
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,

		Status: rules.UserStatusActive,
		Role:   rules.UserRoleUser,

		Metadata: []byte("{}"),

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) UpdateUser(firstName, lastName, phoneNumber, avatarURL string, metadata json.RawMessage) {
	if firstName != "" {
		u.FirstName = firstName
	}
	if lastName != "" {
		u.LastName = lastName
	}
	if phoneNumber != "" {
		u.PhoneNumber = phoneNumber
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	if metadata != nil {
		u.Metadata = metadata
	}

	u.UpdatedAt = time.Now()
}
