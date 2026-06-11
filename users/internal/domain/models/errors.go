package models

import "fmt"

type ErrorCode string

const (
	ErrCodeInternal          ErrorCode = "INTERNAL_ERROR"
	ErrCodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidStatus     ErrorCode = "INVALID_STATUS"
	ErrCodeInvalidRole       ErrorCode = "INVALID_ROLE"
	ErrCodeInvalidRequest    ErrorCode = "INVALID_REQUEST"
)

type Error struct {
	Code    ErrorCode
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}
