package unit

import (
	"testing"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

// IsUserValid checks only if email or password is empty, email validation is executing in auth service
func TestIsUserValid(t *testing.T) {
	t.Parallel()

	invalidUsers := []models.User{
		{Email: "", Password: "", Role: models.ADMIN},
		{Email: "", Password: "", Role: models.EMPLOYEE},
		{Email: "", Password: "", Role: models.EMPLOYER},
		{Email: "1", Password: "", Role: models.ADMIN},
		{Email: "", Password: "1", Role: models.EMPLOYEE},
	}

	for _, user := range invalidUsers {
		assert.False(t, user.IsValid())
	}

	validUsers := []models.User{
		{Email: "test1", Password: "password1", Role: models.ADMIN},
		{Email: "username", Password: "qe", Role: models.EMPLOYEE},
		{Email: "valid", Password: "psw", Role: models.EMPLOYER},
	}

	for _, user := range validUsers {
		assert.True(t, user.IsValid())
	}
}
