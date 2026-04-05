package auth_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestRefreshTokens(t *testing.T) {
	t.Parallel()

	// create valid and invalid jwt tokens for testing
	validTokens, invalidTokens := createTokens()
	auth := newAuthService()

	t.Run("Refresh valid tokens", func(t *testing.T) {
		testRefreshValidTokens(t, auth, validTokens)
	})

	t.Run("Refresh invalid tokens", func(t *testing.T) {
		testRefreshInvalidTokens(t, auth, invalidTokens)
	})
}

func testRefreshValidTokens(t *testing.T, auth *service.AuthService, tokens []string) {
	t.Parallel()

	for i, refresh := range tokens {
		access, ref, err := auth.RefreshTokens(t.Context(), refresh)
		assert.Nil(t, err)

		userID := testTokenValidatation(t, access, ref)
		assert.Equal(t, fixtures.UserIDs[i], userID)
	}
}

func testRefreshInvalidTokens(t *testing.T, auth *service.AuthService, tokens []string) {
	t.Parallel()

	for _, refresh := range tokens {
		access, ref, err := auth.RefreshTokens(t.Context(), refresh)
		assert.NotNil(t, err)

		var e models.Error
		if errors.As(err, &e) {
			assert.Equal(t, models.ErrCodeInvalidJWT, e.Code)
		} else {
			t.Fatal("should be models Error in auth service refresh tokins")
		}

		assert.Empty(t, access)
		assert.Empty(t, ref)
	}
}
