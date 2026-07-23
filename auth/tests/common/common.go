package common

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenValidatation Check if access and refresh tokens are valid and have same payload.
func TestTokenValidatation(t *testing.T, access string, refresh string, status bool) uuid.UUID {
	pAccess := validateTokenPayload(t, access, status)
	pRefresh := validateTokenPayload(t, refresh, status)

	comparePayloads(t, pAccess, pRefresh)

	id, err := uuid.Parse(pAccess.UserID)
	require.NoError(t, err)

	return id
}

func validateTokenPayload(t *testing.T, token string, status bool) *models.JwtPayload {
	t.Helper()

	payload, err := mocks.Jwter.GetPayload(token)

	require.NoError(t, err)
	assert.NotEmpty(t, payload.UserID)
	assert.Equal(t, status, payload.Verificated)

	return payload
}

func comparePayloads(t *testing.T, pAccess *models.JwtPayload, pRefresh *models.JwtPayload) {
	t.Helper()

	assert.Equal(t, pAccess.UserID, pRefresh.UserID) // access and refresh user id must be the same.
	assert.Equal(t, pAccess.Verificated, pRefresh.Verificated)
	assert.Equal(t, pAccess.Role, pRefresh.Role)
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
