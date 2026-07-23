package handlers_test

import (
	"context"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	queue := make(chan *models.EmailMessage, len(fixtures.Users))
	handlers := newHandlers(queue)

	tokens := getVerificationTokens(handlers, queue)

	t.Run("Test verify user", func(t *testing.T) {
		t.Parallel()
		testVerify(t, handlers, tokens)
	})

	t.Run("Test verify user with invalid token", func(t *testing.T) {
		t.Parallel()
		testVerifyInvalidToken(t, handlers)
	})
}

// Test to verify users.
func testVerify(t *testing.T, handlers *handlers.Handler, tokens map[string]string) {
	t.Helper()
	for _, token := range tokens {
		resp, err := handlers.VerifyEmail(t.Context(), &auth_api.VerifToken{
			Token: token,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
		common.TestTokenValidatation(t, resp.GetAccess(), resp.GetRefresh(), true)
	}
}

// Test to verify users with invalid tokens.
func testVerifyInvalidToken(t *testing.T, handlers *handlers.Handler) {
	t.Helper()
	resp, err := handlers.VerifyEmail(t.Context(), &auth_api.VerifToken{
		Token: fixtures.InvalidVerificationToken,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.NotFound, "").Error())

	assert.Nil(t, resp)
}

// Get verification tokens from verification queue.
func getVerificationTokens(handlers *handlers.Handler, queue chan *models.EmailMessage) map[string]string {
	tokens := make(map[string]string, len(fixtures.Users))
	for _, req := range fixtures.Users {
		handlers.Register(context.TODO(), req)

		message := <-queue
		tokens[req.GetEmail()] = message.Token
	}

	return tokens
}
