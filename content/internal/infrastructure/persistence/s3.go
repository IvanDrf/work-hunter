package persistence

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/IvanDrf/work-hunter/content/internal/domain/models"
	"github.com/minio/minio-go/v7"
)

type s3Repository struct {
	client *minio.Client
	bucket string
	log    *slog.Logger
}

func NewS3Repository(client *minio.Client, bucket string, log *slog.Logger) *s3Repository {
	return &s3Repository{
		client: client,
		bucket: bucket,
		log:    log,
	}
}

func (r *s3Repository) getPath(userID string, cType models.ContentType) string {
	return fmt.Sprintf("%s/%s", cType, userID)
}

func (r *s3Repository) Upload(ctx context.Context, file io.Reader, metadata *models.ContentMetadata) error {
	path := r.getPath(metadata.UserID, metadata.Type)
	opts := minio.PutObjectOptions{ContentType: metadata.ContentType}

	_, err := r.client.PutObject(ctx, r.bucket, path, file, metadata.Size, opts)
	if err != nil {
		r.log.Error("failed to upload to s3", slog.String("error", err.Error()), slog.String("path", path))
		return models.ErrS3Operation
	}
	return nil
}

func (r *s3Repository) Download(ctx context.Context, userID string, cType models.ContentType) (io.ReadCloser, *models.ContentMetadata, error) {
	path := r.getPath(userID, cType)

	obj, err := r.client.GetObject(ctx, r.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}

	stat, err := obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil, models.ErrContentNotFound
		}
		r.log.Error("failed to stat s3 object", slog.String("error", err.Error()))
		return nil, nil, models.ErrS3Operation
	}

	meta := &models.ContentMetadata{
		UserID:      userID,
		Type:        cType,
		ContentType: stat.ContentType,
		Size:        stat.Size,
	}
	return obj, meta, nil
}

func (r *s3Repository) Delete(ctx context.Context, userID string, cType models.ContentType) error {
	path := r.getPath(userID, cType)
	if err := r.client.RemoveObject(ctx, r.bucket, path, minio.RemoveObjectOptions{}); err != nil {
		r.log.Error("failed to delete from s3", slog.String("error", err.Error()))
		return models.ErrS3Operation
	}

	return nil
}
