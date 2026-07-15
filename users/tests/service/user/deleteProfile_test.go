package user

import (
	"context"
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/tests/service/user/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProfile_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("DeleteUser", ctx, user.ID, false).Return(nil)

	err := svc.DeleteProfile(ctx, user.ID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfile_DeletedUserPermanent(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()
	user.Status = rules.UserStatusDeleted

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("DeleteUser", ctx, user.ID, true).Return(nil)

	err := svc.DeleteProfile(ctx, user.ID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProfile_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()

	err := svc.DeleteProfile(ctx, fixtures.DeleteUserRequestInvalidDTO())

	assert.Error(t, err)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInvalidRequest, modelsErr.Code)

	mockRepo.AssertNotCalled(t, "GetUserByID")
	mockRepo.AssertNotCalled(t, "DeleteUser")
}

func TestDeleteProfile_UserNotFound(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	id := uuid.New()
	repoErr := models.Error{
		Message: "user not found",
		Code:    models.ErrCodeUserNotFound,
	}

	mockRepo.On("GetUserByID", ctx, id).Return(nil, repoErr)

	err := svc.DeleteProfile(ctx, id.String())

	assert.Error(t, err)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertNotCalled(t, "DeleteUser")
}

func TestDeleteProfile_DeleteError(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()
	repoErr := models.Error{
		Message: "failed to delete user",
		Code:    models.ErrCodeInternal,
	}

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("DeleteUser", ctx, user.ID, false).Return(repoErr)

	err := svc.DeleteProfile(ctx, user.ID.String())

	assert.Error(t, err)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}
