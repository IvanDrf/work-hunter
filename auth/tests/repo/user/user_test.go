package user_test

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/postgres"
	"github.com/IvanDrf/work-hunter/auth/tests/repo/user/fixtures"
	"github.com/stretchr/testify/assert"
)

func connect() (*postgres.UserRepo, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewUserRepo(db)
	return repo, mock
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	for _, user := range users {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Email, user.HashedPassword, user.Verificated, user.Role).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.CreateUser(t.Context(), user)
		assert.Nil(t, err)
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	for _, user := range users {
		mock.ExpectExec("DELETE FROM users").
			WithArgs(user.Email).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteUser(t.Context(), user.Email)
		assert.Nil(t, err)
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindUserByEmail(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	for _, user := range users {
		row := mock.NewRows([]string{"user_id", "email", "hashed_password", "verificated", "role"}).
			AddRow(user.ID, user.Email, user.HashedPassword, user.Verificated, user.Role)

		mock.ExpectQuery("SELECT user_id, email, hashed_password, verificated, role FROM users").
			WithArgs(user.Email).WillReturnRows(row)

		u, err := repo.FindUserByEmail(t.Context(), user.Email)

		assert.Nil(t, err)
		assert.Equal(t, user.ID, u.ID)
		assert.Equal(t, user.Email, u.Email)
		assert.Equal(t, user.HashedPassword, u.HashedPassword)
		assert.Equal(t, user.Verificated, u.Verificated)
		assert.Equal(t, user.Role, u.Role)
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindUserByID(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()
	for _, user := range users {
		row := mock.NewRows([]string{"user_id", "email", "hashed_password", "verificated", "role"}).
			AddRow(user.ID, user.Email, user.HashedPassword, user.Verificated, user.Role)

		mock.ExpectQuery("SELECT user_id, email, hashed_password, verificated, role FROM users").
			WithArgs(user.ID).WillReturnRows(row)

		u, err := repo.FindUserByID(t.Context(), user.ID)

		assert.Nil(t, err)
		assert.Equal(t, user.ID, u.ID)
		assert.Equal(t, user.Email, u.Email)
		assert.Equal(t, user.HashedPassword, u.HashedPassword)
		assert.Equal(t, user.Verificated, u.Verificated)
		assert.Equal(t, user.Role, u.Role)

		assert.Nil(t, mock.ExpectationsWereMet())
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestVerifyEmail(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := fixtures.CreateUsers()

	for _, user := range users {
		mock.ExpectExec("UPDATE users").
			WithArgs(user.Email).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.VerifyEmail(t.Context(), user.Email)

		assert.Nil(t, err)
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}
