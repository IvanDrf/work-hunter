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
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	repo := mocks.NewUserRepo()
	auth := service.NewAuthService(repo, mocks.Jwter)

	t.Run("Delete users", func(t *testing.T) {
		testDeleteUsers(auth, repo, t)
	})

	t.Run("Delete unregistred users", func(t *testing.T) {
		testDeleteUnregistredUsers(auth, repo, t)
	})

	t.Run("Delete users with invalid jwt", func(t *testing.T) {
		testDeleteUsersWithInvalidJWT(auth, repo, t)
	})

	t.Run("Delete users with invalid userID", func(t *testing.T) {
		testDeleteUsersWithInvalidUserID(auth, repo, t)
	})

	t.Run("Delete users with invalid password", func(t *testing.T) {
		testDeleteUsersWithInvalidPassword(auth, repo, t)
	})
}

// Create access tokens
//
//	src = map[email][password]
func createAccessTokens(
	jwter jwt.Jwter,
	users map[string]string, ids []string,
	t *testing.T,
) map[string]string {
	// access tokens: email - token
	tokens := make(map[string]string, len(users))

	i := 0
	for email := range users {
		access, _, err := jwter.CreateTokens(&models.JwtPayload{
			UserID:      ids[i],
			Verificated: false,
			Role:        models.EMPLOYEE,
		})
		assert.Nil(t, err)

		tokens[email] = access
		i++
	}

	return tokens
}

// Trying to delete users and should get an error with an errorCode from models.ErrorCode
//
//	tokens = map[email][access]
//	users = map[email][password]
func checkErrorAfterDeleteUser(
	auth *service.AuthService, repo *mocks.UserRepo,
	tokens map[string]string, users map[string]string,
	errorCode models.ErrorCode, t *testing.T,
) {
	storageLen := len(repo.Storage) // storage should not change

	for email, password := range users {
		err := auth.DeleteUser(t.Context(), tokens[email], password)

		assert.NotNil(t, err)
		assert.Equal(t, storageLen, len(repo.Storage))

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
func registerUsers(auth *service.AuthService, users map[string]string, t *testing.T) map[string]string {
	tokens := make(map[string]string, len(users))

	for email, password := range users {
		access, _, err := auth.RegisterUser(t.Context(), email, password, string(models.EMPLOYEE))
		assert.Nil(t, err)
		assert.NotEmpty(t, access)

		tokens[email] = access
	}

	return tokens
}

// Test to delete existing users
func testDeleteUsers(auth *service.AuthService, repo *mocks.UserRepo, t *testing.T) {
	// access tokens for registred users

	tokens := registerUsers(auth, fixtures.Users, t)

	// test to delete users
	for email, password := range fixtures.Users {
		err := auth.DeleteUser(t.Context(), tokens[email], password)
		assert.Nil(t, err)

		// after user has been deleted he should not be in repo
		user, err := repo.FindUserByEmail(t.Context(), email)
		assert.NotNil(t, err)
		assert.Nil(t, user)

	}
}

// Test to delete unregistred users
func testDeleteUnregistredUsers(auth *service.AuthService, repo *mocks.UserRepo, t *testing.T) {
	// create tokens for unregistred users
	tokens := createAccessTokens(mocks.Jwter, fixtures.Unregistered, fixtures.UserIDsString[:], t)

	// trying to delete unregistred users, should be errors with code ErrCodeUserNotFound
	checkErrorAfterDeleteUser(auth, repo, tokens, fixtures.Unregistered, models.ErrCodeUserNotFound, t)
}

// Test to delete user with invalid jwt token
func testDeleteUsersWithInvalidJWT(auth *service.AuthService, repo *mocks.UserRepo, t *testing.T) {
	// create invalid jwt tokens
	tokens := createAccessTokens(mocks.InvalidJwter, fixtures.Users, fixtures.UserIDsString[:], t)

	// trying to delete users with invalid jwt, should be errors with code ErrCodeInvalidJWT
	checkErrorAfterDeleteUser(auth, repo, tokens, fixtures.Users, models.ErrCodeInvalidJWT, t)
}

// Test to delete users with invalid userID in jwt token, userID is not uuid
func testDeleteUsersWithInvalidUserID(auth *service.AuthService, repo *mocks.UserRepo, t *testing.T) {
	const invalidID = "invalid_id"

	// create invalid users ids
	ids := make([]string, 0, len(fixtures.UserIDs))
	for range len(fixtures.UserIDs) {
		ids = append(ids, invalidID)
	}

	// create valid jwt tokens with invalid UserID
	tokens := createAccessTokens(mocks.Jwter, fixtures.Users, ids, t)

	// trying to delete users with valid jwt tokens but with invalid userID in token
	checkErrorAfterDeleteUser(auth, repo, tokens, fixtures.Users, models.ErrCodeInvalidJWT, t)
}

// Test to delete users with invalid password for account
func testDeleteUsersWithInvalidPassword(auth *service.AuthService, repo *mocks.UserRepo, t *testing.T) {
	const invalidPassword = "invalid_password"

	// create users with invalid password
	users := make(map[string]string, len(fixtures.Users))
	for email := range fixtures.Users {
		users[email] = invalidPassword
	}

	// register users
	tokens := registerUsers(auth, fixtures.Users, t)

	// trying to delete users with invalid passwords
	checkErrorAfterDeleteUser(auth, repo, tokens, users, models.ErrCodeInvalidPassword, t)
}
