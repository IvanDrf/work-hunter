package auth_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/common"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenPayload(t *testing.T) {
	// create valid and invalud access jwt
	validTokens, invalidTokens := common.CreateTokens()
	auth := newAuthService()

	t.Run("Test GetTokenPayload from valid tokens", func(t *testing.T) {
		t.Parallel()

		for i, token := range validTokens {
			testValidTokensPayload(t, auth, token, &models.JwtPayload{
				ID:          fixtures.UserIDs[i],
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
	p, err := auth.GetTokenPayload(t.Context(), token)
	assert.Nil(t, err)

	assert.Equal(t, payload.ID, p.ID)
	assert.Equal(t, payload.Verificated, p.Verificated)
}

func testInvalidTokensPayload(t *testing.T, auth *service.AuthService, token string) {
	p, err := auth.GetTokenPayload(t.Context(), token)

	var e models.Error
	if errors.As(err, &e) {
		assert.Equal(t, models.ErrCodeInvalidJWT, e.Code)
	} else {
		t.Fatalf("should be models Error in auth service GetTokenPayload")
	}

	assert.Nil(t, p)
}
