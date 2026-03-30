package mocks

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type userRepo struct {
	storage map[string]*models.User
}

func NewUserRepo() *userRepo {
	return &userRepo{
		storage: map[string]*models.User{},
	}
}

func (u *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	if u.storage == nil {
		return models.Error{
			Message: "storage is closed",
			Code:    models.ErrCodeInternal,
		}
	}

	if _, ok := u.storage[user.Email]; ok {
		return models.Error{
			Message: "user already exists",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}

	u.storage[user.Email] = user
	return nil
}

func (u *userRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if u.storage == nil {
		return nil, models.Error{
			Message: "storage is closed",
			Code:    models.ErrCodeInternal,
		}
	}

	if user, ok := u.storage[email]; ok {
		return user, nil
	}

	return nil, models.Error{
		Message: "user with that email doesn't exists",
		Code:    models.ErrCodeUserNotFound,
	}
}

func (u *userRepo) VerifyEmail(ctx context.Context, email string) error {
	if u.storage == nil {
		return models.Error{
			Message: "storage is closed",
			Code:    models.ErrCodeInternal,
		}
	}

	if user, ok := u.storage[email]; ok {
		user.Verificated = true

		return nil
	}

	return models.Error{
		Message: "user with that email doesn't exists",
		Code:    models.ErrCodeUserNotFound,
	}
}

func (u *userRepo) Close() {
	u.storage = nil
}
