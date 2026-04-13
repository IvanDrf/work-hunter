package mocks

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/google/uuid"
)

type UserRepo struct {
	Storage map[string]*models.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		Storage: map[string]*models.User{},
	}
}

// Returns new filled user repo with registred users from fixtures.Users
func NewFilledUserRepo() *UserRepo {
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

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	if u.Storage == nil {
		return models.Error{
			Message: "storage is closed",
			Code:    models.ErrCodeInternal,
		}
	}

	if _, ok := u.Storage[user.Email]; ok {
		return models.Error{
			Message: "user already exists",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}

	u.Storage[user.Email] = user
	return nil
}

func (u *UserRepo) DeleteUser(ctx context.Context, email string) error {
	old := len(u.Storage)
	delete(u.Storage, email)

	// if storage len doesn't change means there are no user wit given email
	if old == len(u.Storage) {
		return models.Error{
			Message: "can't delete user with given email",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	return nil
}

func (u *UserRepo) FindUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	for _, user := range u.Storage {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, models.Error{
		Message: "can't find user with given userID",
		Code:    models.ErrCodeUserNotFound,
	}
}

func (u *UserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if u.Storage == nil {
		return nil, models.Error{
			Message: "storage is closed",
			Code:    models.ErrCodeInternal,
		}
	}

	if user, ok := u.Storage[email]; ok {
		return user, nil
	}

	return nil, models.Error{
		Message: "user with that email doesn't exists",
		Code:    models.ErrCodeUserNotFound,
	}
}

func (u *UserRepo) ChangeUserPassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	user, err := u.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.HashedPassword = hashedPassword
	return nil
}

func (u *UserRepo) VerifyEmail(ctx context.Context, email string) error {
	if u.Storage == nil {
		return models.Error{
			Message: "storage is closed",
			Code:    models.ErrCodeInternal,
		}
	}

	if user, ok := u.Storage[email]; ok {
		user.Verificated = true

		return nil
	}

	return models.Error{
		Message: "user with that email doesn't exists",
		Code:    models.ErrCodeUserNotFound,
	}
}

func (u *UserRepo) Close() {
	u.Storage = nil
}
