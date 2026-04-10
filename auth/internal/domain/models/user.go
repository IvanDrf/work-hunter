package models

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"user_id"`

	Email          string `json:"email"`
	HashedPassword string `json:"password"`
	Role           Role   `json:"role"`

	Verificated bool `json:"verificated"`
}

func NewUser(email string, password string, role Role) (*User, error) {
	if !rules.IsPasswordCorrect(password) {
		return nil, Error{
			Message: "password is to long or short",
			Code:    ErrCodeInvalidPassword,
		}
	}

	if !rules.IsEmailValid(email) {
		return nil, Error{
			Message: "invalid email",
			Code:    ErrCodeInvalidEmail,
		}
	}

	hashedPassword, err := rules.HashPassword(password)
	if err != nil {
		return nil, Error{
			Message: "can't hash user password",
			Code:    ErrCodeInternal,
		}
	}

	return &User{
		ID:             uuid.New(),
		Email:          email,
		HashedPassword: hashedPassword,
		Role:           role,
		Verificated:    false,
	}, nil
}
