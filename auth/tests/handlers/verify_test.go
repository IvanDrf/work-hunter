package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	queue := make(chan *models.EmailMessage, len(fixtures.Users))
	handlers := newHandlers(queue)

	// Register users with handlers, cuz we need to get verification messages
	messages := make([]*models.EmailMessage, 0, len(fixtures.Users))
	for _, req := range fixtures.Users {
		_, err := handlers.Register(t.Context(), req)
		assert.Nil(t, err)

		messages = append(messages, <-queue)
	}

	t.Run("Test verify", func(t *testing.T) {
		t.Parallel()
		testVerify(t, handlers, messages)
	})

	t.Run("Test verify with invalid token", func(t *testing.T) {
		t.Parallel()
		testVerifyInvalidToken(t, handlers)
	})
}

// Test to verify user account with valid tokens
func testVerify(t *testing.T, handlers *handlers.Handler, messages []*models.EmailMessage) {
	for _, message := range messages {
		resp, err := handlers.VerifyEmail(t.Context(), &auth_api.VerifToken{
			Token: message.Token,
		})
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		common.TestTokenValidatation(t, resp.Access, resp.Refresh, true)
	}
}

// Test to verify user account with invalid token
func testVerifyInvalidToken(t *testing.T, handlers *handlers.Handler) {
	resp, err := handlers.VerifyEmail(t.Context(), &auth_api.VerifToken{
		Token: fixtures.InvalidVerificationToken,
	})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.NotFound, "").Error())

	assert.Nil(t, resp)
}
