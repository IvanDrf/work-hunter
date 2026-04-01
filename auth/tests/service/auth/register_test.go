package auth_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/auth/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	auth := newAuthService()

	t.Run("Register new users", func(t *testing.T) {
		testRegisterNewUsers(t, auth)
	})

	t.Run("Register old users", func(t *testing.T) {
		testRegisterOldUsers(t, auth)
	})

}

func testRegisterNewUsers(t *testing.T, auth *service.AuthService) {
	for email, password := range fixtures.Users {
		access, refresh, err := auth.RegisterUser(t.Context(), email, password)
		assert.Nil(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		// jwt tokens after registration should be valid
		testTokenValidatation(t, access, refresh)
	}
}

func testRegisterOldUsers(t *testing.T, auth *service.AuthService) {
	//users already registred, should be errors
	for email, password := range fixtures.Users {
		access, refresh, err := auth.RegisterUser(t.Context(), email, password)
		assert.NotNil(t, err)

		var e models.Error
		if errors.As(err, &e) {
			assert.Equal(t, models.ErrCodeUserAlreadyExists, e.Code)

		} else {
			t.Fatal("should be models Error in auth service RegisterUser")
		}

		assert.Empty(t, access)
		assert.Empty(t, refresh)
	}
}
