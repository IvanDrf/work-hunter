package models

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"user_id"`

	Username       string `json:"username"`
	HashedPassword string `json:"password"`
}

func NewUser(username string, password string) (*User, error) {
	if !rules.IsPasswordCorrect(password) {
		return nil, Error{
			Message: "password is to long or short",
			Code:    ErrCodeInvalidPassword,
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
		Username:       username,
		HashedPassword: hashedPassword,
	}, nil
}
