package models

import "fmt"

type ErrorCode string

const (
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"

	ErrCodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"

	ErrCodeInvalidPassword ErrorCode = "INVALID_PASSWORD"
)

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message: %s, code: %s", e.Message, e.Code)
}
