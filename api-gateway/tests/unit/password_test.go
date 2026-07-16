package unit

import (
	"testing"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestIsPasswordValid(t *testing.T) {
	t.Parallel()

	invalidPasswords := []models.Password{
		{Old: "", New: ""},
		{Old: "1", New: ""},
		{Old: "", New: "New"},
	}

	for _, password := range invalidPasswords {
		assert.False(t, password.IsValid())
	}

	validPasswords := []models.Password{
		{Old: "eurgjnle", New: "we;llf"},
		{Old: "1", New: "wel,f"},
		{Old: "ef", New: "New"},
	}

	for _, user := range validPasswords {
		assert.True(t, user.IsValid())
	}

}
