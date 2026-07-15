package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()
	updatedUser := &models.User{
		ID:          user.ID,
		FirstName:   "Updated",
		LastName:    "Name",
		CompanyName: "New Corp",
		Verificated: true,
	}

	mock.ExpectExec("UPDATE users SET").
		WithArgs(
			updatedUser.FirstName,
			updatedUser.LastName,
			updatedUser.CompanyName,
			updatedUser.Verificated,
			updatedUser.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateUser(t.Context(), updatedUser)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_NoRowsAffected(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()
	user.FirstName = "Updated"

	mock.ExpectExec("UPDATE users SET").
		WithArgs(
			user.FirstName,
			user.LastName,
			user.CompanyName,
			user.Verificated,
			user.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateUser(t.Context(), user)
	assert.NoError(t, err) // Обновление несуществующего пользователя не вызывает ошибку
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_InternalError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()

	mock.ExpectExec("UPDATE users SET").
		WithArgs(
			user.FirstName,
			user.LastName,
			user.CompanyName,
			user.Verificated,
			user.ID,
		).
		WillReturnError(errors.New("update failed"))

	err := repo.UpdateUser(t.Context(), user)

	assert.Error(t, err)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInternal, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "failed to update user")
	assert.NoError(t, mock.ExpectationsWereMet())
}
