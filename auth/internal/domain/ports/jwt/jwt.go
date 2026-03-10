package jwt

import "github.com/google/uuid"

type Jwter interface {
	CreateTokens(userID uuid.UUID) (string, string, error)
	GetUserID(token string) (uuid.UUID, error)

	RefreshTokens(refresh string) (string, string, error)
}
