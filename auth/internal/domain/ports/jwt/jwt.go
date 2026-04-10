package jwt

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type Jwter interface {
	CreateTokens(payload *models.JwtPayload) (string, string, error)
	GetPayload(token string) (*models.JwtPayload, error)
}
