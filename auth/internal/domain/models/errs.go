package models

import "fmt"

type ErrorCode string

const (
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"

	ErrCodeInvalidUserRole        ErrorCode = "INVALID_USER_ROLE"
	ErrCodeUserNotFound           ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists      ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeUserAlreadyVerificated ErrorCode = "USER_ALREADY_VERIFICATED"
	ErrCodeOutdatedToken          ErrorCode = "TOKEN_IS_OUTDATED"

	ErrCodeInvalidPassword ErrorCode = "INVALID_PASSWORD"
	ErrCodeInvalidEmail    ErrorCode = "INVALID_EMAIL"

	ErrCodeInvalidJWT ErrorCode = "INVALID_JWT_TOKENS"
)

const (
	CantCreateJWTMsg        = "can't create jwt token for user"
	InvalidUserIDInTokenMsg = "invalid userID in jwt token"
	InvalidJWTMsg           = "invalid jwt token"
)

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message: %s, code: %s", e.Message, e.Code)
}
