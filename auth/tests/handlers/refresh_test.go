package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TedstRefresh(t *testing.T) {
	t.Parallel()

	handlers := newHandlers()

	// access tokens
	tokens := registerUsers(handlers, fixtures.Users)

	t.Run("Test refresh tokens", func(t *testing.T) {
		t.Parallel()
		testRefresh(t, handlers, tokens)
	})

	t.Run("Test refresh tokens with invalid signature", func(t *testing.T) {
		t.Parallel()
		testRefreshInvalidTokens(t, handlers, tokens)
	})
}

// Test to refresh valid jwt tokens
func testRefresh(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	for _, tokens := range tokens {
		resp, err := handlers.RefreshTokens(t.Context(), &auth_api.RefreshToken{
			Refresh: tokens.Refresh,
		})

		assert.Nil(t, err)
		assert.NotNil(t, resp)

		common.TestTokenValidatation(t, resp.Access, resp.Refresh, false)
	}
}

// Test to refresh invalid jwt tokens
func testRefreshInvalidTokens(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	for _, tokens := range tokens {
		// get payload from valid refresh token
		payload, _ := mocks.Jwter.GetPayload(tokens.Refresh)

		// create refresh token with invalid signature, but with valid payload
		_, refresh, _ := mocks.InvalidJwter.CreateTokens(payload)

		resp, err := handlers.RefreshTokens(t.Context(), &auth_api.RefreshToken{
			Refresh: refresh,
		})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

		assert.Nil(t, resp)
	}
}
