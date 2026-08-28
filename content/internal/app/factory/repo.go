package factory

import (
	"fmt"

	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/content/internal/infrastructure/persistence"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (f *Factory) newRepos() (repo.ContentRepo, error) {
	return f.newS3Repo()
}

func (f *Factory) newS3Repo() (*persistence.S3Repository, error) {
	minioClient, err := minio.New(f.cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(f.cfg.RootUser, f.cfg.RootPassword, ""),
		Secure: f.cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init minio: %w", err)
	}

	repo := persistence.NewS3Repository(minioClient, f.cfg.BucketName, f.log)

	return repo, nil
}
