package user

import (
	"context"
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/tests/service/user/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)

	result, err := svc.GetProfile(ctx, user.ID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID.String(), result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.FirstName, result.FirstName)
	assert.Equal(t, user.LastName, result.LastName)
	assert.Equal(t, user.CompanyName, result.CompanyName)
	assert.Equal(t, string(user.Status), result.Status)
	assert.Equal(t, string(user.Role), result.Role)
	assert.Equal(t, user.Verificated, result.Verificated)

	mockRepo.AssertExpectations(t)
}

func TestGetProfile_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()

	result, err := svc.GetProfile(ctx, fixtures.GetUserRequestInvalidDTO())

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInvalidRequest, modelsErr.Code)

	mockRepo.AssertNotCalled(t, "GetUserByID")
}

func TestGetProfile_NotFound(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	id := uuid.New()
	repoErr := models.Error{
		Message: "user not found",
		Code:    models.ErrCodeUserNotFound,
	}

	mockRepo.On("GetUserByID", ctx, id).Return(nil, repoErr)

	result, err := svc.GetProfile(ctx, id.String())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestGetProfile_InternalError(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	id := uuid.New()
	repoErr := models.Error{
		Message: "database error",
		Code:    models.ErrCodeInternal,
	}

	mockRepo.On("GetUserByID", ctx, id).Return(nil, repoErr)

	result, err := svc.GetProfile(ctx, id.String())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}
