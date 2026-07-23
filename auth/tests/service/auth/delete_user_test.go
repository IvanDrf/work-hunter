package auth_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUser(t *testing.T) {
	repo := mocks.NewUserRepo()
	auth := service.NewAuthService(repo, mocks.Jwter)

	t.Run("Delete users", func(t *testing.T) {
		testDeleteUsers(t, auth, repo)
	})

	t.Run("Delete unregistred users", func(t *testing.T) {
		testDeleteUnregistredUsers(t, auth, repo)
	})

	t.Run("Delete users with invalid jwt", func(t *testing.T) {
		testDeleteUsersWithInvalidJWT(t, auth, repo)
	})

	t.Run("Delete users with invalid userID", func(t *testing.T) {
		testDeleteUsersWithInvalidUserID(t, auth, repo)
	})

	t.Run("Delete users with invalid password", func(t *testing.T) {
		testDeleteUsersWithInvalidPassword(t, auth, repo)
	})
}

// Create access tokens
//
//	src = map[email][password]
func createAccessTokens(
	t *testing.T,
	jwter jwt.Jwter,
	users map[string]string, ids []string,
) map[string]string {
	t.Helper()

	// access tokens: email - token.
	tokens := make(map[string]string, len(users))

	i := 0
	for email := range users {
		access, _, err := jwter.CreateTokens(&models.JwtPayload{
			UserID:      ids[i],
			Verificated: false,
			Role:        models.EMPLOYEE,
		})
		assert.NoError(t, err)

		tokens[email] = access
		i++
	}

	return tokens
}

// Trying to delete users and should get an error with an errorCode from models.ErrorCode.
//
//	tokens = map[email][access]
//	users = map[email][password]
func checkErrorAfterDeleteUser(
	t *testing.T,
	auth *service.AuthService, repo *mocks.UserRepo,
	tokens map[string]string, users map[string]string,
	errorCode models.ErrorCode,
) {
	t.Helper()
	storageLen := len(repo.Storage) // storage should not change.

	for email, password := range users {
		err := auth.DeleteUser(t.Context(), tokens[email], password)

		require.Error(t, err)
		assert.Len(t, repo.Storage, storageLen)

		var e models.Error
		if errors.As(err, &e) {
			assert.Equal(t, errorCode, e.Code)
		} else {
			t.Fatalf("DeleteUser error code should be %s, but: %s", errorCode, err.Error())
		}
	}
}

// Register users with auth service
//
//	users = map[email][password]
func registerUsers(t *testing.T, auth *service.AuthService, users map[string]string) map[string]string {
	t.Helper()

	tokens := make(map[string]string, len(users))

	for email, password := range users {
		access, _, err := auth.RegisterUser(t.Context(), email, password, string(models.EMPLOYEE))
		require.NoError(t, err)
		assert.NotEmpty(t, access)

		tokens[email] = access
	}

	return tokens
}

// Test to delete existing users.
func testDeleteUsers(t *testing.T, auth *service.AuthService, repo *mocks.UserRepo) {
	t.Helper()

	// access tokens for registred users.

	tokens := registerUsers(t, auth, fixtures.Users)

	// test to delete users.
	for email, password := range fixtures.Users {
		err := auth.DeleteUser(t.Context(), tokens[email], password)
		require.NoError(t, err)

		// after user has been deleted he should not be in repo.
		user, err := repo.FindUserByEmail(t.Context(), email)
		require.Error(t, err)
		assert.Nil(t, user)
	}
}

// Test to delete unregistred users.
func testDeleteUnregistredUsers(t *testing.T, auth *service.AuthService, repo *mocks.UserRepo) {
	t.Helper()

	// create tokens for unregistred users.
	tokens := createAccessTokens(t, mocks.Jwter, fixtures.Unregistered, fixtures.UserIDsString[:])

	// trying to delete unregistred users, should be errors with code ErrCodeUserNotFound.
	checkErrorAfterDeleteUser(t, auth, repo, tokens, fixtures.Unregistered, models.ErrCodeUserNotFound)
}

// Test to delete user with invalid jwt token.
func testDeleteUsersWithInvalidJWT(t *testing.T, auth *service.AuthService, repo *mocks.UserRepo) {
	t.Helper()

	// create invalid jwt tokens.
	tokens := createAccessTokens(t, mocks.InvalidJwter, fixtures.Users, fixtures.UserIDsString[:])

	// trying to delete users with invalid jwt, should be errors with code ErrCodeInvalidJWT.
	checkErrorAfterDeleteUser(t, auth, repo, tokens, fixtures.Users, models.ErrCodeInvalidJWT)
}

// Test to delete users with invalid userID in jwt token, userID is not uuid.
func testDeleteUsersWithInvalidUserID(t *testing.T, auth *service.AuthService, repo *mocks.UserRepo) {
	t.Helper()

	const invalidID = "invalid_id"

	// create invalid users ids.
	ids := make([]string, 0, len(fixtures.UserIDs))
	for range len(fixtures.UserIDs) {
		ids = append(ids, invalidID)
	}

	// create valid jwt tokens with invalid UserID.
	tokens := createAccessTokens(t, mocks.Jwter, fixtures.Users, ids)

	// trying to delete users with valid jwt tokens but with invalid userID in token.
	checkErrorAfterDeleteUser(t, auth, repo, tokens, fixtures.Users, models.ErrCodeInvalidJWT)
}

// Test to delete users with invalid password for account.
func testDeleteUsersWithInvalidPassword(t *testing.T, auth *service.AuthService, repo *mocks.UserRepo) {
	t.Helper()

	const invalidPassword = "invalid_password"

	// create users with invalid password.
	users := make(map[string]string, len(fixtures.Users))
	for email := range fixtures.Users {
		users[email] = invalidPassword
	}

	// register users.
	tokens := registerUsers(t, auth, fixtures.Users)

	// trying to delete users with invalid passwords.
	checkErrorAfterDeleteUser(t, auth, repo, tokens, users, models.ErrCodeInvalidPassword)
}
