package user

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail_Success(t *testing.T) {
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
			WithArgs(user.Email).
			WillReturnRows(rows)

		result, err := repo.GetUserByEmail(t.Context(), user.Email)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	email := "nonexistent@example.com"

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetUserByEmail(t.Context(), email)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeUserNotFound, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_DeletedUser(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()
	user.Status = rules.UserStatusDeleted

	mock.ExpectQuery("SELECT \\* FROM users").
		WithArgs(user.Email).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetUserByEmail(t.Context(), user.Email)

	assert.Error(t, err)
	assert.Nil(t, result)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeUserNotFound, modelsErr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
