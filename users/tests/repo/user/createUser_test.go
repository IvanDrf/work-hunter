package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	for _, user := range users {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				user.ID, user.Email,
				user.FirstName, user.LastName, user.CompanyName,
				user.Status, user.Role, user.Verificated,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateUser(t.Context(), user)
		assert.NoError(t, err)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()

	pqErr := &pq.Error{
		Code:    "23505",
		Message: "duplicate key value violates unique constraint",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.ID, user.Email,
			user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated,
		).
		WillReturnError(pqErr)

	err := repo.CreateUser(t.Context(), user)

	assert.Error(t, err)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeUserAlreadyExists, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "already exists")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_InternalError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()

	mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.ID, user.Email,
			user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated,
		).
		WillReturnError(errors.New("database connection failed"))

	err := repo.CreateUser(t.Context(), user)

	assert.Error(t, err)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInternal, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "failed to create user")
	assert.NoError(t, mock.ExpectationsWereMet())
}
