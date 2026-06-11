package jwt_test

import (
	"testing"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	secret        = "test_secret_key"
	invalid_token = "invalid_jwt"
)

var (
	jwter                   = jwt.NewJwt(secret, 15*time.Minute, 20*time.Minute)
	userID                  = uuid.New()
	verificated             = false
	role        models.Role = "ADMIN"
)

func TestCreateTokens(t *testing.T) {
	t.Parallel()

	access, refresh, err := jwter.CreateTokens(&models.JwtPayload{
		UserID:      userID.String(),
		Verificated: verificated,
		Role:        role,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
}

func TestGetTokenClaims(t *testing.T) {
	t.Parallel()

	access, refresh, err := jwter.CreateTokens(&models.JwtPayload{
		UserID:      userID.String(),
		Verificated: verificated,
		Role:        role,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	checkClaims(t, access)
	checkClaims(t, refresh)

	// check invalid token claims, payload should be nil
	payload, err := jwter.GetPayload(invalid_token)
	assert.NotNil(t, err)
	assert.Nil(t, payload)
}

func checkClaims(t *testing.T, token string) {
	t.Helper()

	payload, err := jwter.GetPayload(token)
	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, userID.String(), payload.UserID)
	assert.Equal(t, verificated, payload.Verificated)
	assert.Equal(t, role, payload.Role)
}
