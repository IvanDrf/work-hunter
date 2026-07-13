package models

import "fmt"

type ErrCode string

const (
	ErrCodeInternal             ErrCode = "INTERNAL_ERROR"
	ErrCodeUnsupportedMediaType ErrCode = "UNSUPPORTED_MEDIA"
)

type Error struct {
	Message string  `json:"message"`
	Code    ErrCode `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("message:%s, code:%s", e.Message, e.Code)
}
