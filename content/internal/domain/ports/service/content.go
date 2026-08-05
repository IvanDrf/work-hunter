package service

import (
	"context"
	"io"

	"github.com/IvanDrf/work-hunter/content/internal/domain/models"
)

type ContentService interface {
	UploadResume(ctx context.Context, userID string, file io.Reader, size int64) error
	UploadAvatar(ctx context.Context, userID string, file io.Reader, size int64) error
	GetContent(ctx context.Context, userID string, cType models.ContentType) (io.ReadCloser, *models.ContentMetadata, error)
	DeleteContent(ctx context.Context, userID string, cType models.ContentType) error
}
