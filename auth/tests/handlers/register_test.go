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

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	handlers := newHandlers()
	t.Run("Register users", func(t *testing.T) {
		testRegisterNewUsers(t, handlers) // using
	})

	t.Run("Register already registred users", func(t *testing.T) {
		// using the same handlers as Register users
		//  cuz we need to check can we register already registred users
		testRegisterOldUsers(t, handlers)
	})

	t.Run("Register users with invalid role", func(t *testing.T) {
		t.Parallel()
		testRegisterInvalidRoleUsers(t, newHandlers())
	})

	t.Run("Register users with invalid password", func(t *testing.T) {
		t.Parallel()
		testRegisterInvalidPasswordUsers(t, newHandlers())
	})

	t.Run("Register users with invalid email", func(t *testing.T) {
		t.Parallel()
		testRegisterInvalidEmailUsers(t, newHandlers())
	})
}

// Test to register new users
func testRegisterNewUsers(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.RegisterRequests {
		resp, err := handlers.Register(t.Context(), req)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		common.TestTokenValidatation(t, resp.Access, resp.Refresh, false)
	}
}

// Test to register users what already registred
func testRegisterOldUsers(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.RegisterRequests {
		resp, err := handlers.Register(t.Context(), req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.AlreadyExists, "").Error())
	}
}

// Test to register users with invalid roles
func testRegisterInvalidRoleUsers(t *testing.T, handlers *handlers.Handler) {
	for _, req := range fixtures.InvalidRoleRequests {
		resp, err := handlers.Register(t.Context(), req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
	}
}

// Test to register users with invalid passwords
func testRegisterInvalidPasswordUsers(t *testing.T, handlers *handlers.Handler) {
	getErrorWithInvalidArgument(t, handlers, fixtures.InvalidPasswordRequests)
}

// Test to register users with invalid email
func testRegisterInvalidEmailUsers(t *testing.T, handlers *handlers.Handler) {
	getErrorWithInvalidArgument(t, handlers, fixtures.InvalidEmailRequests)
}

// Trying to register user using content from requests and expecting error status code
//
//	codes.InvalidArgument
func getErrorWithInvalidArgument(t *testing.T, handlers *handlers.Handler, requests []*auth_api.User) {
	for _, req := range requests {
		resp, err := handlers.Register(t.Context(), req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
	}
}
