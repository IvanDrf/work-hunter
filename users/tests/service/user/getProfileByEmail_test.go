package user

import (
	"context"
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/tests/service/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestGetProfileByEmail_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()

	mockRepo.On("GetUserByEmail", ctx, user.Email).Return(user, nil)

	result, err := svc.GetProfileByEmail(ctx, user.Email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID.String(), result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.LastName, result.LastName)
	assert.Equal(t, user.CompanyName, result.CompanyName)

	mockRepo.AssertExpectations(t)
}

func TestGetProfileByEmail_NotFound(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	email := fixtures.GetUserByEmailNotFoundDTO()
	repoErr := models.Error{
		Message: "user not found",
		Code:    models.ErrCodeUserNotFound,
	}

	mockRepo.On("GetUserByEmail", ctx, email).Return(nil, repoErr)

	result, err := svc.GetProfileByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestGetProfileByEmail_InternalError(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	email := "test@example.com"
	repoErr := models.Error{
		Message: "database error",
		Code:    models.ErrCodeInternal,
	}

	mockRepo.On("GetUserByEmail", ctx, email).Return(nil, repoErr)

	result, err := svc.GetProfileByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}
