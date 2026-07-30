package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("Test delete user", func(t *testing.T) {
		t.Parallel()

		handlers := newHandlers(nil)
		tokens := registerUsers(handlers, fixtures.Users)

		testDeleteUser(t, handlers, tokens)
	})

	t.Run("Test delete user with invalid jwt", func(t *testing.T) {
		t.Parallel()

		handlers := newHandlers(nil)
		tokens := registerUsers(handlers, fixtures.Users)

		testDeleteUserInvalidJWT(t, handlers, tokens)
	})

	t.Run("Test to delete user with invalid userID in jwt", func(t *testing.T) {
		t.Parallel()

		testDeleteUserInvalidJWTUserID(t, newHandlers(nil))
	})

	t.Run("Test to delete user with not existing userID", func(t *testing.T) {
		t.Parallel()

		testDeleteUserInvalidUserID(t, newHandlers(nil))
	})

	t.Run("Test to delete user with invalid password", func(t *testing.T) {
		t.Parallel()

		handlers := newHandlers(nil)
		tokens := registerUsers(handlers, fixtures.Users)

		testDeleteUserInvalidPassword(t, handlers, tokens)
	})
}

// Test to delete users.
func testDeleteUser(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	t.Helper()

	for _, req := range fixtures.Users {
		resp, err := handlers.DeleteUser(t.Context(), &auth_api.DeleteUserRequest{
			Access:   tokens[req.GetEmail()].Access,
			Password: req.GetPassword(),
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)

		// after deletion user can't login.
		_, err = handlers.Login(t.Context(), req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.NotFound, "").Error())
	}
}

// Test to delete user with invalid jwt.
func testDeleteUserInvalidJWT(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	t.Helper()

	for _, req := range fixtures.Users {
		// get userID from valid access token.
		payload, _ := handlers.IsTokenValid(t.Context(), &auth_api.AccessToken{
			Access: tokens[req.GetEmail()].Access,
		})

		// create access jwt token with valid payload, but with invalid signature.
		access, _, _ := mocks.InvalidJwter.CreateTokens(&models.JwtPayload{
			UserID:      payload.GetId(),
			Verificated: false,
			Role:        models.EMPLOYEE,
		})

		resp, err := handlers.DeleteUser(t.Context(), &auth_api.DeleteUserRequest{
			Access:   access,
			Password: req.GetPassword(),
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

		assert.Nil(t, resp)
	}
}

// Test to delete user valid jwt but invalid userID.
func testDeleteUserInvalidJWTUserID(t *testing.T, handlers *handlers.Handler) {
	t.Helper()

	for _, req := range fixtures.Users {
		// create valid jwt tokens with invalid userID.
		access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
			UserID:      fixtures.InvalidUserID,
			Verificated: false,
			Role:        models.EMPLOYEE,
		})

		resp, err := handlers.DeleteUser(t.Context(), &auth_api.DeleteUserRequest{
			Access:   access,
			Password: req.GetPassword(),
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

		assert.Nil(t, resp)
	}
}

// Test to delete user with valid userID, but user with that userID doesn't exists.
func testDeleteUserInvalidUserID(t *testing.T, handlers *handlers.Handler) {
	t.Helper()

	userID := uuid.New()

	for _, req := range fixtures.Users {
		access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
			UserID:      userID.String(),
			Verificated: false,
			Role:        models.EMPLOYEE,
		})

		resp, err := handlers.DeleteUser(t.Context(), &auth_api.DeleteUserRequest{
			Access:   access,
			Password: req.GetPassword(),
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.NotFound, "").Error())

		assert.Nil(t, resp)
	}
}

// Test to delete user with invalid password.
func testDeleteUserInvalidPassword(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	t.Helper()

	for _, req := range fixtures.Users {
		resp, err := handlers.DeleteUser(t.Context(), &auth_api.DeleteUserRequest{
			Access:   tokens[req.GetEmail()].Access,
			Password: fixtures.InvalidPassword,
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())

		assert.Nil(t, resp)
	}
}
