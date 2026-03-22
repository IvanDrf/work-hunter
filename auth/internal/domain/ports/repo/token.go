package repo

import (
	"context"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type TokenRepo interface {
	CreateToken(ctx context.Context, email string, token *models.Token) error
	FindEmailExpByToken(ctx context.Context, token string) (string, time.Time, error)

	DeleteToken(ctx context.Context, token string) error
}
