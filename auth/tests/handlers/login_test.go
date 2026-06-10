package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginHandler(t *testing.T) {
	t.Parallel()

	handlers := newHandlers(nil)
	// register users, so we can log in
	registerUsers(handlers, fixtures.Users)

	t.Run("Login registred users", func(t *testing.T) {
		testLoginUsers(t, handlers)
	})

	t.Run("Login unregistred users", func(t *testing.T) {
		testLoginUnregistredUsers(t, handlers)
	})

	t.Run("Login with invalid password", func(t *testing.T) {
		testLoginWithInvalidPassword(t, handlers)
	})

}

// Test to login registred users
func testLoginUsers(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.Users {
		resp, err := handlers.Login(t.Context(), req)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		// jwt tokens must be valid after successfull login
		common.TestTokenValidatation(t, resp.Access, resp.Refresh, false)
	}
}

// Test to login unregistred users
func testLoginUnregistredUsers(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.UnregistredUsers {
		resp, err := handlers.Login(t.Context(), req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.NotFound, "").Error())
	}
}

// Test to login user with invalid password
func testLoginWithInvalidPassword(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.Users {
		r := &auth_api.User{
			Email:    req.Email,
			Password: fixtures.InvalidPassword,
			Role:     req.Role,
		}

		resp, err := handlers.Login(t.Context(), r)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
	}
}
