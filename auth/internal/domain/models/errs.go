package models

import "fmt"

type ErrorCode string

const (
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"

	ErrCodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrOutdatedToken         ErrorCode = "TOKEN_IS_OUTDATED"

	ErrCodeInvalidPassword ErrorCode = "INVALID_PASSWORD"
	ErrCodeInvalidEmail    ErrorCode = "INVALID_EMAIL"

	ErrCodeInvalidJWT ErrorCode = "INVALID_JWT_TOKENS"
)

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message: %s, code: %s", e.Message, e.Code)
}
