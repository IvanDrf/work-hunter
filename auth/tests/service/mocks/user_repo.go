package mocks

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
)

type userRepo struct {
	storage map[string]*models.User
}

func NewUserRepo() *userRepo {
	return &userRepo{
		storage: map[string]*models.User{},
	}
}

// Returns new filled user repo with registred users from fixtures.Users
func NewFilledUserRepo() *userRepo {
	userRepo := NewUserRepo()
	i := 0
	for email, password := range fixtures.Users {
		userRepo.CreateUser(context.TODO(), &models.User{
			ID:             fixtures.UserIDs[i],
			Email:          email,
			HashedPassword: password,
			Verificated:    false,
		})

		i++
	}

	return userRepo
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
