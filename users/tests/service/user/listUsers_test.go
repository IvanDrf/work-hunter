package user

import (
	"context"
	"errors"
	"testing"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/tests/service/user/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListUsers_Success(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestDTO()
	users := fixtures.TestUsers()
	expectedTotal := int32(len(users))

	mockRepo.On("ListUsers", ctx, mock.MatchedBy(func(params *models.ListUsersParams) bool {
		return params.Page == req.Page &&
			params.PageSize == req.PageSize &&
			params.Status == req.Filter.Status &&
			params.Role == req.Filter.Role &&
			params.SearchQuery == req.Filter.SearchQuery &&
			params.SortBy == req.SortBy &&
			params.SortDesc == req.SortDesc
	})).Return(users, expectedTotal, nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, len(users))
	assert.Equal(t, expectedTotal, result.TotalCount)
	assert.Equal(t, req.Page, result.Page)
	assert.Equal(t, req.PageSize, result.PageSize)
	assert.Equal(t, int32(1), result.TotalPages)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_WithoutFilter(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestWithoutFilterDTO()
	users := fixtures.TestUsers()
	expectedTotal := int32(len(users))

	mockRepo.On("ListUsers", ctx, mock.MatchedBy(func(params *models.ListUsersParams) bool {
		return params.Status == "" &&
			params.Role == "" &&
			params.SearchQuery == ""
	})).Return(users, expectedTotal, nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, len(users))
	assert.Equal(t, expectedTotal, result.TotalCount)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_WithSearchQuery(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestWithSearchDTO()
	user := fixtures.TestUser()
	users := []*models.User{user}
	expectedTotal := int32(1)

	mockRepo.On("ListUsers", ctx, mock.MatchedBy(func(params *models.ListUsersParams) bool {
		return params.SearchQuery == "john"
	})).Return(users, expectedTotal, nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "John", result.Users[0].FirstName)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_WithStatusFilter(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestWithStatusDTO()
	users := []*models.User{fixtures.TestUser()}
	expectedTotal := int32(1)

	mockRepo.On("ListUsers", ctx, mock.MatchedBy(func(params *models.ListUsersParams) bool {
		return params.Status == string(rules.UserStatusActive)
	})).Return(users, expectedTotal, nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, string(rules.UserStatusActive), result.Users[0].Status)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_WithRoleFilter(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestWithRoleDTO()
	users := []*models.User{fixtures.TestUserEmployer()}
	expectedTotal := int32(1)

	mockRepo.On("ListUsers", ctx, mock.MatchedBy(func(params *models.ListUsersParams) bool {
		return params.Role == string(rules.UserRoleEmployer)
	})).Return(users, expectedTotal, nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, string(rules.UserRoleEmployer), result.Users[0].Role)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_MaxPageSize(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestMaxPageSizeDTO()
	users := fixtures.TestUsers()
	expectedTotal := int32(len(users))

	mockRepo.On("ListUsers", ctx, mock.MatchedBy(func(params *models.ListUsersParams) bool {
		return params.PageSize == 100
	})).Return(users, expectedTotal, nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 100, result.PageSize)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_RepositoryError(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestDTO()
	repoErr := errors.New("database error")

	mockRepo.On("ListUsers", ctx, mock.AnythingOfType("*models.ListUsersParams")).Return(nil, int32(0), repoErr)

	result, err := svc.ListUsers(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_EmptyResult(t *testing.T) {
	t.Parallel()

	svc, mockRepo := setupTest(t)
	ctx := context.Background()
	req := fixtures.ListUsersRequestDTO()

	mockRepo.On("ListUsers", ctx, mock.AnythingOfType("*models.ListUsersParams")).Return([]*models.User{}, int32(0), nil)

	result, err := svc.ListUsers(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Users)
	assert.Equal(t, int32(0), result.TotalCount)
	assert.Equal(t, int32(0), result.TotalPages)

	mockRepo.AssertExpectations(t)
}
