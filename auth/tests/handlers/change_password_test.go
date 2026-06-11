package handlers_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestChangePasswordHandler(t *testing.T) {
	t.Parallel()

	handlers := newHandlers(nil)
	// register users, so we can change their passwords
	tokens := registerUsers(handlers, fixtures.Users)

	t.Run("Change user password", func(t *testing.T) {
		testChangePassword(t, handlers, tokens)
	})

	t.Run("Change user password with invalid old password", func(t *testing.T) {
		testChangePasswordInvalidOld(t, handlers, tokens)
	})

	t.Run("Change user password with invalid new password", func(t *testing.T) {
		testChangePasswordInvalidNew(t, handlers, tokens)
	})

	t.Run("Change user password with invalid jwt token", func(t *testing.T) {
		testChangePasswordInvalidJwt(t, handlers)
	})

	t.Run("Change user password with invalid userID in jwt token", func(t *testing.T) {
		testChangePasswordInvalidID(t, handlers)
	})
}

// Test to change password for registred users
func testChangePassword(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	// change password from old to new
	for _, req := range fixtures.Users {
		resp, err := handlers.ChangePassword(t.Context(), &auth_api.ChangePasswordRequest{
			Access: tokens[req.Email].Access,
			Old:    req.Password,
			New:    fixtures.NewPassword,
		})

		assert.Nil(t, err)
		assert.NotNil(t, resp)
	}

	// trying to login with new password
	for _, req := range fixtures.Users {
		resp, err := handlers.Login(t.Context(), &auth_api.User{
			Email:    req.Email,
			Password: fixtures.NewPassword,
			Role:     req.Role,
		})

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		common.TestTokenValidatation(t, resp.Access, resp.Refresh, false)
	}
}

// Test to change password with invalid old
func testChangePasswordInvalidOld(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	for _, req := range fixtures.Users {
		resp, err := handlers.ChangePassword(t.Context(), &auth_api.ChangePasswordRequest{
			Access: tokens[req.Email].Access,
			Old:    fixtures.InvalidOldPassword,
			New:    req.Password,
		})

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
	}
}

// Test to change password with invalid new, means new password doesn't fit rules.IsPasswordCorrect()
func testChangePasswordInvalidNew(t *testing.T, handlers *handlers.Handler, tokens map[string]*Tokens) {
	for _, req := range fixtures.Users {
		resp, err := handlers.ChangePassword(t.Context(), &auth_api.ChangePasswordRequest{
			Access: tokens[req.Email].Access,
			Old:    fixtures.NewPassword, // use new password because change it in first test
			New:    fixtures.InvalidNewPassword,
		})

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
	}
}

// Test to change password with invalid jwt token
func testChangePasswordInvalidJwt(t *testing.T, handlers *handlers.Handler) {
	access, _, _ := mocks.InvalidJwter.CreateTokens(&models.JwtPayload{
		UserID:      "user_id",
		Verificated: false,
		Role:        models.EMPLOYEE,
	})

	resp, err := handlers.ChangePassword(t.Context(), &auth_api.ChangePasswordRequest{
		Access: access,
		Old:    fixtures.NewPassword,
		New:    fixtures.NewPassword,
	})

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
}

// Test to change password with valid jwt token but with invalid userID - not uuid
func testChangePasswordInvalidID(t *testing.T, handlers *handlers.Handler) {
	access, _, _ := mocks.Jwter.CreateTokens(&models.JwtPayload{
		UserID:      fixtures.InvalidUserID,
		Verificated: false,
		Role:        models.EMPLOYEE,
	})

	resp, err := handlers.ChangePassword(t.Context(), &auth_api.ChangePasswordRequest{
		Access: access,
		Old:    fixtures.NewPassword,
		New:    fixtures.NewPassword,
	})

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), status.Error(codes.InvalidArgument, "").Error())
}
