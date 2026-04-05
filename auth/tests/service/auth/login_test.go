package auth_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/common"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	t.Parallel()

	auth := newAuthService()

	t.Run("Login users", func(t *testing.T) {
		testLoginUsers(t, auth)
	})

	t.Run("Login unregistred users", func(t *testing.T) {
		testLoginUnregistredUsers(t, auth)
	})

}

func testLoginUsers(t *testing.T, auth *service.AuthService) {
	for email, password := range fixtures.Users {
		// register new user
		_, _, err := auth.RegisterUser(t.Context(), email, password)
		assert.Nil(t, err)

		// login this new user
		access, refresh, err := auth.LoginUser(t.Context(), email, password)
		assert.Nil(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		// jwt tokens after login should be valid
		common.TestTokenValidatation(t, access, refresh, false)
	}
}

func testLoginUnregistredUsers(t *testing.T, auth *service.AuthService) {
	for email, password := range fixtures.Unregistered {
		access, refres, err := auth.LoginUser(t.Context(), email, password)
		assert.NotNil(t, err)

		var e models.Error
		if errors.As(err, &e) {
			assert.Equal(t, models.ErrCodeUserNotFound, e.Code)
		} else {
			t.Fatal("should be models Error in auth service login user")
		}

		assert.Empty(t, access)
		assert.Empty(t, refres)
	}
}
