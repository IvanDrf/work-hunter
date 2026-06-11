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
	// should not be parallel, because we use common Queue for SendVerificationEmail and VerifyEmailByToken

	verif := newVerificationService()

	// send verif email to registred users, should be no errors
	t.Run("Send verification email", func(t *testing.T) {
		testSendVerificationEmail(t, verif)
	})

	// send verif email to unregistred users, should be errors
	t.Run("Send verification email", func(t *testing.T) {
		testSendVerificationEmailUnregistredUsers(t, verif)
	})

	verif.Close()
}

func testSendVerificationEmail(t *testing.T, verif *service.VerificationService) {
	sended := 0 // amount of sended emails, should be equal to email producer queue len
	for email := range fixtures.Users {
		sended++

		err := verif.SendVerificationEmail(t.Context(), email)
		assert.Nil(t, err)
		assert.Equal(t, sended, len(Queue)) // len of email producer queue should increase and be equal to sended
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

		assert.Equal(t, sended, len(Queue)) // len of email producer queue should not change and should be equal to sended
	}
}
