package models

import "fmt"

type ErrCode string

const (
	ErrCodeInternal             ErrCode = "INTERNAL_ERROR"
	ErrCodeUnsupportedMediaType ErrCode = "UNSUPPORTED_MEDIA"
	ErrCodeUnprocessableEntity  ErrCode = "UNPROCESSABLE_ENTITY"
	ErrCodeAlreadyExists        ErrCode = "ALREADY_EXISTS"
	ErrCodeInvalidArgument      ErrCode = "INVALID_ARGUMENT"
	ErrCodeInvalidCookie        ErrCode = "INVALID_COOKIE"
	ErrCodeAccess               ErrCode = "ACCESS_ERROR"
	ErrCodeNotFound             ErrCode = "NOT_FOUND"
	ErrCodeAborted              ErrCode = "ABORTED"
	ErrCodeFailedPrecondition   ErrCode = "FAILED_PRECONDITION"
	ErrCodeResourceExhausted    ErrCode = "RESOURCE_EXHAUSTED"
	ErrCodeDeadlineExceeded     ErrCode = "DEADLINE_EXCEEDED"
	ErrCodeCanceled             ErrCode = "CANCELED"
)

type Error struct {
	Message string  `json:"message"`
	Code    ErrCode `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message:%s, code:%s", e.Message, e.Code)
}
