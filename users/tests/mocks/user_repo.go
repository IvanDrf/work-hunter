package mocks

import (
	"context"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) DeleteUser(ctx context.Context, id uuid.UUID, permanent bool) error {
	args := m.Called(ctx, id, permanent)
	return args.Error(0)
}

func (m *UserRepositoryMock) ListUsers(ctx context.Context, params *models.ListUsersParams) ([]*models.User, int32, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.User), args.Get(1).(int32), args.Error(2)
}

func (m *UserRepositoryMock) UpdateUserStatus(ctx context.Context, id uuid.UUID, status rules.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *UserRepositoryMock) Close() {
	m.Called()
}
