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
	jwter       = pkg.NewJwt(secret, 15*time.Minute, 20*time.Minute)
	userID      = uuid.New()
	verificated = false
)

func TestCreateTokens(t *testing.T) {
	t.Parallel()

	access, refresh, err := jwter.CreateTokens(userID, verificated)
	assert.Nil(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
}

func TestGetTokenClaims(t *testing.T) {
	t.Parallel()

	access, refresh, err := jwter.CreateTokens(userID, verificated)
	assert.Nil(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	checkClaims(t, access)
	checkClaims(t, refresh)

	userIDFromJwt, ver, err := jwter.GetPayload(invalid_token)
	assert.NotNil(t, err)
	assert.Empty(t, userIDFromJwt)
	assert.False(t, ver)
}

func checkClaims(t *testing.T, token string) {
	t.Helper()

	userIDFromJwt, ver, err := jwter.GetPayload(token)
	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, userID, userIDFromJwt)
	assert.Equal(t, verificated, ver)
}
