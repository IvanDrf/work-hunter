package common

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Check if access and refresh tokens are valid and have same payload
func TestTokenValidatation(t *testing.T, access string, refresh string, status bool) uuid.UUID {
	accessID, verificated, err := mocks.Jwter.GetPayload(access)
	assert.Nil(t, err)
	assert.NotNil(t, accessID)
	assert.Equal(t, status, verificated)

	refreshID, verificated, err := mocks.Jwter.GetPayload(refresh)
	assert.Nil(t, err)
	assert.NotNil(t, refreshID)
	assert.Equal(t, status, verificated)

	assert.Equal(t, accessID, refreshID)

	return accessID
}

func CreateTokens() ([]string, []string) {
	validTokens := make([]string, 0, len(fixtures.UserIDs))
	invalidTokens := make([]string, 0, len(fixtures.UserIDs))

	for _, userID := range fixtures.UserIDs {
		_, validRefresh, _ := mocks.Jwter.CreateTokens(userID, false)
		_, invalidRefresh, _ := mocks.InvalidJwter.CreateTokens(userID, false)

		validTokens = append(validTokens, validRefresh)
		invalidTokens = append(invalidTokens, invalidRefresh)
	}

	return validTokens, invalidTokens
}
