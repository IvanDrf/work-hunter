package jwt

import "github.com/google/uuid"

type Jwter interface {
	CreateTokens(userID uuid.UUID, verificated bool) (string, string, error)
	GetPayload(token string) (uuid.UUID, bool, error)

	RefreshTokens(refresh string) (string, string, error)
}
