package verification_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/common"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestVerifyEmailByToken(t *testing.T) {
	// should not be parallel, because we use common Queue for SendVerificationEmail and VerifyEmailByToken

	verif := newVerificationService()

	// verify emails for registred users
	t.Run("Verify email by token for users", func(t *testing.T) {
		testVerifyEmailByToken(t, verif)
	})

	t.Run("Verify email by token with invalid tokens", func(t *testing.T) {
		testVerifyEmailByTokenInvalidToken(t, verif)
	})

}

func testVerifyEmailByToken(t *testing.T, verif *service.VerificationService) {
	for _, token := range fixtures.Tokens {
		access, refresh, err := verif.VerifyEmailByToken(t.Context(), token)
		assert.Nil(t, err)

		common.TestTokenValidatation(t, access, refresh, true) // jwt tokens should be valud and status = true verificated
	}
}

func testVerifyEmailByTokenInvalidToken(t *testing.T, verif *service.VerificationService) {
	for _, token := range fixtures.UnusedTokens {
		access, refresh, err := verif.VerifyEmailByToken(t.Context(), token)

		var e models.Error
		if errors.As(err, &e) {
			assert.Equal(t, models.ErrCodeUserNotFound, e.Code)
		} else {
			t.Fatalf("verify email by token: should me models error with code %s, actual %s", models.ErrCodeUserNotFound, err.Error())
		}

		assert.Empty(t, access)
		assert.Empty(t, refresh)
	}
}
