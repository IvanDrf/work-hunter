package testutils

import (
	"testing"
	"time"

	"github.com/IvanDrf/work-hunter/auth/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	secret        = "test_secret_key"
	invalid_token = "invalid_jwt"
)

var (
	jwter  = pkg.NewJwt(secret, 15*time.Minute, 20*time.Minute)
	userID = uuid.New()
)

func TestCreateTokens(t *testing.T) {
	t.Parallel()

	access, refresh, err := jwter.CreateTokens(userID)
	assert.Nil(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
}

func TestGetTokenClaims(t *testing.T) {
	t.Parallel()

	access, refresh, err := jwter.CreateTokens(userID)
	assert.Nil(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	checkClaims(t, access)
	checkClaims(t, refresh)

	claims, err := jwter.GetTokenClaims(invalid_token)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
}

func checkClaims(t *testing.T, token string) {
	t.Helper()

	claims, err := jwter.GetTokenClaims(token)
	assert.Nil(t, err)
	assert.NotNil(t, claims)

	userIDClaims, err := uuid.Parse(claims.UserID)
	assert.Nil(t, err)
	assert.Equal(t, userID, userIDClaims)
}
