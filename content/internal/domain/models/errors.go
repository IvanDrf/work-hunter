package models

import "errors"

var (
	ErrFileEmpty       = errors.New("file is empty")
	ErrFileTooLarge    = errors.New("file size exceeds allowed limit")
	ErrReadFailed      = errors.New("failed to read file header")
	ErrInvalidMimeType = errors.New("invalid file content type")

	ErrS3Operation     = errors.New("storage operation failed")
	ErrContentNotFound = errors.New("content not found")
)
