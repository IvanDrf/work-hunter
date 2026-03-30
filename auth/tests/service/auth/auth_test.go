package auth_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/auth/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/service/auth/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newAuthService() *service.AuthService {
	return service.NewAuthService(mocks.NewUserRepo(), fixtures.Jwter)
}

// Check if access and refresh tokens are valid and have same payload
func testTokenValidatation(t *testing.T, access string, refresh string) uuid.UUID {
	accessID, verificated, err := fixtures.Jwter.GetPayload(access)
	assert.Nil(t, err)
	assert.NotNil(t, accessID)
	assert.False(t, verificated)

	refreshID, verificated, err := fixtures.Jwter.GetPayload(refresh)
	assert.Nil(t, err)
	assert.NotNil(t, refreshID)
	assert.False(t, verificated)

	assert.Equal(t, accessID, refreshID)

	return accessID
}

func createTokens() ([]string, []string) {
	validTokens := make([]string, 0, len(fixtures.UserIDs))
	invalidTokens := make([]string, 0, len(fixtures.UserIDs))

	for _, userID := range fixtures.UserIDs {
		_, validRefresh, _ := fixtures.Jwter.CreateTokens(userID, false)
		_, invalidRefresh, _ := fixtures.InvalidJwter.CreateTokens(userID, false)

		validTokens = append(validTokens, validRefresh)
		invalidTokens = append(invalidTokens, invalidRefresh)
	}

	return validTokens, invalidTokens
}
