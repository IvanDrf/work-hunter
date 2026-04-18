package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSendEmail(t *testing.T) {
	t.Parallel()

	queue := make(chan *models.EmailMessage, len(fixtures.Users))
	handlers := newHandlers(queue)

	tokens := registerUsers(handlers, fixtures.Users)
	clearQueue(queue, len(queue))

	t.Run("Test send email", func(t *testing.T) {
		t.Parallel()
		testSendEmail(t, handlers, tokens, queue)
	})

	t.Run("Test send email with invalid jwt token", func(t *testing.T) {
		t.Parallel()
		testSendEmailInvalidJWT(t, handlers)
	})

	t.Run("Test send email with invalid userID in jwt token", func(t *testing.T) {
		t.Parallel()
		testSendEmailInvalidUserID(t, handlers)
	})

	t.Run("Test send email with status verificated in jwt", func(t *testing.T) {
		t.Parallel()

		testSendEmailVerificatedInJWT(t, handlers)
	})

	t.Run("Test send email to already verificated users", func(t *testing.T) {
		t.Parallel()

		testSendEmailAlreadyVerificated(t, newHandlersWithVerifUsers(nil))
	})
}

// Test to resend email
func testSendEmail(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens, queue chan *models.EmailMessage) {
	for email, token := range tokens {
		resp, err := handlers.SendVerificationEmail(t.Context(), &auth_api.AccessToken{
			Access: token.Access,
		})

		assert.Nil(t, err)
		assert.NotNil(t, resp)

		message := <-queue
		assert.Equal(t, email, message.Email)
	}
}

// Test to resend email with invalid jwt token
func testSendEmailInvalidJWT(t *testing.T, handlers *handlers.Handler) {
	access, _, _ := mocks.InvalidJwter.CreateTokens(&models.JwtPayload{
		UserID:      "user_id",
		Verificated: false,
		Role:        models.EMPLOYEE,
	})

	resp, err := handlers.SendVerificationEmail(t.Context(), &auth_api.AccessToken{
		Access: access,
	})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

	assert.Nil(t, resp)

}

// Test to send email with invalid userID in jwt token
func testSendEmailInvalidUserID(t *testing.T, handlers *handlers.Handler) {
	access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
		UserID:      fixtures.InvalidUserID,
		Verificated: false,
		Role:        models.EMPLOYEE,
	})

	resp, err := handlers.SendVerificationEmail(t.Context(), &auth_api.AccessToken{
		Access: access,
	})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

	assert.Nil(t, resp)
}

// Test to send email with verificated status in jwt token
// Should be no errors, cuz check verificated status in database
func testSendEmailVerificatedInJWT(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.Users {
		tokens, err := handlers.Login(t.Context(), req)
		assert.Nil(t, err)

		payload, err := mocks.Jwter.GetPayload(tokens.Access)
		assert.Nil(t, err)

		access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
			UserID:      payload.UserID,
			Verificated: true,
			Role:        payload.Role,
		})

		resp, err := handlers.SendVerificationEmail(t.Context(), &auth_api.AccessToken{
			Access: access,
		})

		assert.Nil(t, err)
		assert.NotNil(t, resp)
	}
}

// Test to send verification email already verificated users
func testSendEmailAlreadyVerificated(t *testing.T, handlers *handlers.Handler) {
	// get access jwt tokens
	for _, req := range fixtures.Users {
		tokens, _ := handlers.Login(t.Context(), req)

		resp, err := handlers.SendVerificationEmail(t.Context(), &auth_api.AccessToken{
			Access: tokens.Access,
		})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.AlreadyExists, "").Error())

		assert.Nil(t, resp)
	}
}

func clearQueue(queue chan *models.EmailMessage, n int) {
	for range n {
		<-queue
	}
}
