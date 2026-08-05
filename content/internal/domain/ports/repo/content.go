package repo

import (
	"context"
	"io"

	"github.com/IvanDrf/work-hunter/content/internal/domain/models"
)

type ContentRepo interface {
	Upload(ctx context.Context, file io.Reader, metadata *models.ContentMetadata) error
	Download(ctx context.Context, userID string, cType models.ContentType) (io.ReadCloser, *models.ContentMetadata, error)
	Delete(ctx context.Context, userID string, cType models.ContentType) error
}
