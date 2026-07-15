package user

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	for _, user := range users {
		rows := sqlmock.NewRows([]string{
			"id", "email", "first_name", "last_name", "company_name",
			"status", "role", "verificated", "created_at", "updated_at",
		}).AddRow(
			user.ID, user.Email, user.FirstName, user.LastName, user.CompanyName,
			user.Status, user.Role, user.Verificated, user.CreatedAt, user.UpdatedAt,
		)

		mock.ExpectQuery("SELECT \\* FROM users").
			WithArgs(user.ID).
			WillReturnRows(rows)

		result, err := repo.GetUserByID(t.Context(), user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.FirstName, result.FirstName)
		assert.Equal(t, user.LastName, result.LastName)
		assert.Equal(t, user.CompanyName, result.CompanyName)
		assert.Equal(t, user.Status, result.Status)
		assert.Equal(t, user.Role, result.Role)
		assert.Equal(t, user.Verificated, result.Verificated)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	id := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetUserByID(t.Context(), id)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeUserNotFound, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_InternalError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	id := uuid.New()

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs(id).
		WillReturnError(errors.New("database connection lost"))

	result, err := repo.GetUserByID(t.Context(), id)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInternal, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "failed to get user")
	assert.NoError(t, mock.ExpectationsWereMet())
}
