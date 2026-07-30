package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsTokenValid(t *testing.T) {
	t.Parallel()

	handlers := newHandlers(nil)
	tokens := registerUsers(handlers, fixtures.Users)

	t.Run("Test is token valid", func(t *testing.T) {
		t.Parallel()
		testIsTokenValid(t, handlers, tokens)
	})

	t.Run("Test is token valid with invalid userID", func(t *testing.T) {
		t.Parallel()
		testIsTokenValidInvalidUserID(t, handlers)
	})

	t.Run("Test is tokenv valid with invalid user role", func(t *testing.T) {
		t.Parallel()
		testIsTokenValidInvalidUserRole(t, handlers, tokens)
	})
}

// Test to validate jwt tokens.
func testIsTokenValid(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	t.Helper()

	for _, tokens := range tokens {
		resp, err := handlers.IsTokenValid(t.Context(), &auth_api.AccessToken{
			Access: tokens.Access,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)

		payload, _ := mocks.Jwter.GetPayload(tokens.Access)
		assert.Equal(t, payload.UserID, resp.GetId())
		assert.Equal(t, payload.Verificated, resp.GetVerificated())
		assert.Equal(t, string(payload.Role), resp.GetRole().String())
	}
}

// Test to validate token with invalid userID.
func testIsTokenValidInvalidUserID(t *testing.T, handlers *handlers.Handler) {
	t.Helper()

	access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
		UserID:      fixtures.InvalidUserID,
		Verificated: false,
		Role:        models.EMPLOYEE,
	})

	resp, err := handlers.IsTokenValid(t.Context(), &auth_api.AccessToken{
		Access: access,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

	assert.Nil(t, resp)
}

// Test to validate jwt token with invalid Role.
func testIsTokenValidInvalidUserRole(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	t.Helper()

	for _, tokens := range tokens {
		payload, _ := mocks.Jwter.GetPayload(tokens.Access)

		access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
			UserID:      payload.UserID,
			Verificated: false,
			Role:        fixtures.InvalidUserRole,
		})

		resp, err := handlers.IsTokenValid(t.Context(), &auth_api.AccessToken{
			Access: access,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

		assert.Nil(t, resp)
	}
}
