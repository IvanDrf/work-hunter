package services

import (
	"context"
	"io"
	"log/slog"

	"github.com/IvanDrf/work-hunter/content/internal/domain/models"
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/content/internal/domain/rules"
)

type contentService struct {
	repo repo.ContentRepo
	log  *slog.Logger
}

func NewContentService(repo repo.ContentRepo, log *slog.Logger) *contentService {
	return &contentService{
		repo: repo,
		log:  log,
	}
}

func (s *contentService) UploadResume(ctx context.Context, userID string, file io.Reader, size int64) error {
	return s.uploadWithValidation(ctx, userID, models.TypeResume, file, size)
}

func (s *contentService) UploadAvatar(ctx context.Context, userID string, file io.Reader, size int64) error {
	return s.uploadWithValidation(ctx, userID, models.TypeAvatar, file, size)
}

func (s *contentService) uploadWithValidation(ctx context.Context, userID string, cType models.ContentType, file io.Reader, size int64) error {
	validated, err := rules.ValidateContent(cType, file, size)
	if err != nil {
		s.log.Warn("content validation failed",
			slog.String("userID", userID),
			slog.String("type", string(cType)),
			slog.Int64("size", size),
			slog.String("error", err.Error()),
		)
		return err
	}

	s.log.Info("uploading validated content",
		slog.String("userID", userID),
		slog.String("type", string(cType)),
		slog.String("mime", validated.MimeType),
		slog.Int64("size", validated.Size),
	)

	return s.repo.Upload(ctx, validated.Reader, &models.ContentMetadata{
		UserID:   userID,
		Type:     cType,
		MimeType: validated.MimeType,
		Size:     validated.Size,
	})
}

func (s *contentService) GetContent(ctx context.Context, userID string, cType models.ContentType) (io.ReadCloser, *models.ContentMetadata, error) {
	return s.repo.Download(ctx, userID, cType)
}

func (s *contentService) DeleteContent(ctx context.Context, userID string, cType models.ContentType) error {
	s.log.Info("deleting content", slog.String("userID", userID), slog.String("type", string(cType)))
	return s.repo.Delete(ctx, userID, cType)
}

func (s *contentService) DeleteAllUserContent(ctx context.Context, userID string) error {
	s.log.Info("cascade deleting all content for deleted user", slog.String("userID", userID))

	if err := s.repo.Delete(ctx, userID, models.TypeResume); err != nil {
		s.log.Error("failed to delete resume during cascade wipe", slog.String("userID", userID), slog.String("error", err.Error()))
		return err
	}

	if err := s.repo.Delete(ctx, userID, models.TypeAvatar); err != nil {
		s.log.Error("failed to delete avatar during cascade wipe", slog.String("userID", userID), slog.String("error", err.Error()))
		return err
	}

	return nil
}
