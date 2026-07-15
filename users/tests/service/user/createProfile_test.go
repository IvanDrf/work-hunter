package user

import (
	"context"
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"github.com/IvanDrf/work-hunter/users/tests/service/user/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfile_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.CreateUserRequestDTO()

	mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.Email == req.Email &&
			user.FirstName == req.FirstName &&
			user.LastName == req.LastName &&
			user.CompanyName == req.CompanyName &&
			user.Verificated == req.Verificated
	})).Return(nil)

	result, err := svc.CreateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.ID, result.ID)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.FirstName, result.FirstName)
	assert.Equal(t, req.LastName, result.LastName)
	assert.Equal(t, req.CompanyName, result.CompanyName)
	assert.Equal(t, req.Verificated, result.Verificated)
	assert.Equal(t, string(rules.UserStatusActive), result.Status)
	assert.Equal(t, string(rules.UserRoleEmployer), result.Role)

	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_WithoutCompany_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.CreateUserRequestWithoutCompanyDTO()

	mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.Email == req.Email &&
			user.FirstName == req.FirstName &&
			user.LastName == req.LastName &&
			user.CompanyName == req.CompanyName &&
			user.Role == rules.UserRoleEmployee &&
			user.Verificated == req.Verificated
	})).Return(nil)

	result, err := svc.CreateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.ID, result.ID)
	assert.Equal(t, req.Email, result.Email)
	assert.Equal(t, req.FirstName, result.FirstName)
	assert.Equal(t, req.LastName, result.LastName)
	assert.Equal(t, "", result.CompanyName)
	assert.Equal(t, string(rules.UserRoleEmployee), result.Role)

	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := &dto.CreateUserRequest{
		ID: "invalid-uuid",
	}

	result, err := svc.CreateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInvalidRequest, modelsErr.Code)

	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestCreateProfile_EmptyFields(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.CreateUserRequestEmptyFieldsDTO()

	mockRepo.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
		return user.FirstName == "" &&
			user.LastName == "" &&
			user.CompanyName == ""
	})).Return(nil)

	result, err := svc.CreateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.FirstName)
	assert.Equal(t, "", result.LastName)
	assert.Equal(t, "", result.CompanyName)

	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_RepositoryError(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.CreateUserRequestDTO()
	repoErr := models.Error{
		Message: "database connection failed",
		Code:    models.ErrCodeInternal,
	}

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(repoErr)

	result, err := svc.CreateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateProfile_DuplicateEmail(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.CreateUserRequestDTO()
	repoErr := models.Error{
		Message: "user already exists",
		Code:    models.ErrCodeUserAlreadyExists,
	}

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(repoErr)

	result, err := svc.CreateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeUserAlreadyExists, modelsErr.Code)

	mockRepo.AssertExpectations(t)
}
