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

	userIDFromJwt, err := jwter.GetUserID(invalid_token)
	assert.NotNil(t, err)
	assert.Empty(t, userIDFromJwt)
}

func checkClaims(t *testing.T, token string) {
	t.Helper()

	userIDFromJwt, err := jwter.GetUserID(token)
	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, userID, userIDFromJwt)
}
