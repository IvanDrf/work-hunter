package verification_test

import (
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestSendVerificationEmail(t *testing.T) {
	t.Parallel()

	verif := newVerificationService()

	// send verif email to registred users
	t.Run("Send verification email", func(t *testing.T) {
		testSendVerificationEmail(t, verif)
	})

	t.Run("Send verification email", func(t *testing.T) {
		testSendVerificationEmailUnregistredUsers(t, verif)
	})
}

func testSendVerificationEmail(t *testing.T, verif *service.VerificationService) {
	sended := 0 // amount of sended emails, should be equal to email producer queue len
	for email := range fixtures.Users {
		sended++

		err := verif.SendVerificationEmail(t.Context(), email)
		assert.Nil(t, err)
		assert.Equal(t, len(Queue), sended)
	}
}

func testSendVerificationEmailUnregistredUsers(t *testing.T, verif *service.VerificationService) {
	sended := len(Queue)
	for email := range fixtures.Unregistered {

		err := verif.SendVerificationEmail(t.Context(), email)
		var e models.Error
		if errors.As(err, &e) {
			assert.Equal(t, models.ErrCodeUserNotFound, e.Code)
		} else {
			t.Fatalf("send verif email: should me models error with code %s, actual %s", models.ErrCodeUserNotFound, err.Error())
		}

		assert.Equal(t, sended, len(Queue))
	}
}
