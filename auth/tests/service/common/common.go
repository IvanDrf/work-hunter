package common

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Check if access and refresh tokens are valid and have same payload
func TestTokenValidatation(t *testing.T, access string, refresh string, status bool) uuid.UUID {
	payload, err := mocks.Jwter.GetPayload(access)
	assert.Nil(t, err)
	assert.NotEmpty(t, payload.UserID)
	assert.Equal(t, status, payload.Verificated)

	accessID := payload.UserID

	payload, err = mocks.Jwter.GetPayload(refresh)
	assert.Nil(t, err)
	assert.NotEmpty(t, payload.UserID)
	assert.Equal(t, status, payload.Verificated)

	assert.Equal(t, accessID, payload.UserID) // access and refresh user id must be the same

	id, err := uuid.Parse(accessID)
	assert.Nil(t, err)

	return id
}

func CreateTokens() ([]string, []string) {
	validTokens := make([]string, 0, len(fixtures.UserIDs))
	invalidTokens := make([]string, 0, len(fixtures.UserIDs))

	for _, userID := range fixtures.UserIDs {
		payload := &models.JwtPayload{
			UserID:      userID.String(),
			Verificated: false,
			Role:        models.EMPLOYEE,
		}

		_, validRefresh, _ := mocks.Jwter.CreateTokens(payload)
		_, invalidRefresh, _ := mocks.InvalidJwter.CreateTokens(payload)

		validTokens = append(validTokens, validRefresh)
		invalidTokens = append(invalidTokens, invalidRefresh)
	}

	return validTokens, invalidTokens
}
