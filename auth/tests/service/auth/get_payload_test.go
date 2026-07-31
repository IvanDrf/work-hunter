package auth_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokenPayload(t *testing.T) {
	t.Parallel()

	// create valid and invalud access jwt
	validTokens, invalidTokens := common.CreateTokens()
	auth := newAuthService()

	t.Run("Test GetTokenPayload from valid tokens", func(t *testing.T) {
		t.Parallel()

		for i, token := range validTokens {
			testValidTokensPayload(t, auth, token, &models.JwtPayload{
				UserID:      fixtures.UserIDs[i].String(),
				Verificated: false,
			})
		}
	})

	t.Run("Test GetTokenPayload from invalid tokens", func(t *testing.T) {
		t.Parallel()

		for _, token := range invalidTokens {
			testInvalidTokensPayload(t, auth, token)
		}
	})
}

func testValidTokensPayload(t *testing.T, auth *service.AuthService, token string, payload *models.JwtPayload) {
	t.Helper()

	p, err := auth.GetTokenPayload(t.Context(), token)
	require.NoError(t, err)

	assert.Equal(t, payload.UserID, p.UserID)
	assert.Equal(t, payload.Verificated, p.Verificated)
}

func testInvalidTokensPayload(t *testing.T, auth *service.AuthService, token string) {
	t.Helper()

	p, err := auth.GetTokenPayload(t.Context(), token)

	var e models.Error
	if errors.As(err, &e) {
		assert.Equal(t, models.ErrCodeInvalidJWT, e.Code)
	} else {
		t.Fatalf("should be models Error in auth service GetTokenPayload")
	}

	assert.Nil(t, p)
}
