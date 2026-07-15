package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/users/tests/repo/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserStatus_Success(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()
	newStatus := rules.UserStatusBlocked

	mock.ExpectExec("UPDATE users SET").
		WithArgs(newStatus, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateUserStatus(t.Context(), user.ID, newStatus)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserStatus_InternalError(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	user := fixtures.CreateSingleUser()
	newStatus := rules.UserStatusBlocked

	mock.ExpectExec("UPDATE users SET").
		WithArgs(newStatus, user.ID).
		WillReturnError(errors.New("update status failed"))

	err := repo.UpdateUserStatus(t.Context(), user.ID, newStatus)

	assert.Error(t, err)
	modelsErr, ok := err.(models.Error)
	assert.True(t, ok)
	assert.Equal(t, models.ErrCodeInternal, modelsErr.Code)
	assert.Contains(t, modelsErr.Message, "failed to update user status")
	assert.NoError(t, mock.ExpectationsWereMet())
}
