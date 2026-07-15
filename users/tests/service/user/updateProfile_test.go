package user

import (
	"context"
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"github.com/IvanDrf/work-hunter/users/tests/service/user/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateProfile_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()            // Verificated = true
	req := fixtures.UpdateUserRequestDTO() // Verificated = false

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == user.ID &&
			u.FirstName == req.FirstName &&
			u.LastName == req.LastName &&
			u.CompanyName == req.CompanyName &&
			u.Verificated == user.Verificated
	})).Return(nil)

	result, err := svc.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.FirstName, result.FirstName)
	assert.Equal(t, req.LastName, result.LastName)
	assert.Equal(t, req.CompanyName, result.CompanyName)
	assert.Equal(t, user.Verificated, result.Verificated)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_PartialUpdate(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()
	req := fixtures.UpdateUserRequestPartialDTO()

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == user.ID &&
			u.FirstName == req.FirstName &&
			u.LastName == user.LastName &&
			u.CompanyName == user.CompanyName &&
			u.Verificated == user.Verificated
	})).Return(nil)

	result, err := svc.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.FirstName, result.FirstName)
	assert.Equal(t, user.LastName, result.LastName)
	assert.Equal(t, user.CompanyName, result.CompanyName)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_OnlyCompany(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUserEmployee() // Был Employee
	req := fixtures.UpdateUserRequestOnlyCompanyDTO()

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == user.ID &&
			u.FirstName == user.FirstName &&
			u.LastName == user.LastName &&
			u.CompanyName == req.CompanyName &&
			u.Verificated == user.Verificated
	})).Return(nil)

	result, err := svc.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.CompanyName, result.CompanyName)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.UpdateUserRequestInvalidDTO()

	result, err := svc.UpdateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInvalidRequest, modelsErr.Code)

	mockRepo.AssertNotCalled(t, "GetUserByID")
	mockRepo.AssertNotCalled(t, "UpdateUser")
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.UpdateUserRequestDTO()
	id := uuid.MustParse(req.ID)
	repoErr := models.Error{
		Message: "user not found",
		Code:    models.ErrCodeUserNotFound,
	}

	mockRepo.On("GetUserByID", ctx, id).Return(nil, repoErr)

	result, err := svc.UpdateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertNotCalled(t, "UpdateUser")
}

func TestUpdateProfile_UpdateError(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()
	req := fixtures.UpdateUserRequestDTO()
	repoErr := models.Error{
		Message: "failed to update user",
		Code:    models.ErrCodeInternal,
	}

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("UpdateUser", ctx, mock.AnythingOfType("*models.User")).Return(repoErr)

	result, err := svc.UpdateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_VerificatedFromFalseToTrue(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUserEmployee() // Verificated = true
	user.Verificated = false

	verificated := true
	req := &dto.UpdateUserRequest{
		ID:          user.ID.String(),
		FirstName:   "",
		LastName:    "",
		CompanyName: "",
		Verificated: verificated,
	}

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == user.ID &&
			u.Verificated == true
	})).Return(nil)

	result, err := svc.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Verificated)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProfile_VerificatedTrueToFalse(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	user := fixtures.TestUser()

	verificated := false
	req := &dto.UpdateUserRequest{
		ID:          user.ID.String(),
		FirstName:   "",
		LastName:    "",
		CompanyName: "",
		Verificated: verificated,
	}

	mockRepo.On("GetUserByID", ctx, user.ID).Return(user, nil)
	mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.ID == user.ID &&
			u.Verificated == true
	})).Return(nil)

	result, err := svc.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Verificated)

	mockRepo.AssertExpectations(t)
}
