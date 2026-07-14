package models

import "fmt"

type ErrCode string

const (
	ErrCodeInternal             ErrCode = "INTERNAL_ERROR"
	ErrCodeUnsupportedMediaType ErrCode = "UNSUPPORTED_MEDIA"
	ErrCodeUnprocessableEntity  ErrCode = "UNPROCESSABLE_ENTITY"
	ErrCodeAlreadyExists        ErrCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidArgument      ErrCode = "INVALID_ARGUMENT"
	ErrCodeInvalidCookie        ErrCode = "INVALID_COOKIE"
)

type Error struct {
	Message string  `json:"message"`
	Code    ErrCode `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message:%s, code:%s", e.Message, e.Code)
}
