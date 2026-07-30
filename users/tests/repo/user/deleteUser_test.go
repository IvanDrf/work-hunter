package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_SoftDelete_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()

	mock.ExpectExec("UPDATE users SET status = 'deleted'").
		WithArgs(user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(t.Context(), user.ID, false)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_PermanentDelete_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()

	mock.ExpectExec("DELETE FROM users").
		WithArgs(user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(t.Context(), user.ID, true)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_InternalError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()

	mock.ExpectExec("UPDATE users SET status = 'deleted'").
		WithArgs(user.ID).
		WillReturnError(errors.New("delete failed"))

	err := repo.DeleteUser(t.Context(), user.ID, false)

	assert.Error(t, err)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInternal, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "failed to delete user")
	assert.NoError(t, mock.ExpectationsWereMet())
}
