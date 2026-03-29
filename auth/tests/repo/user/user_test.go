package user_test

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/postgres"
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

const usersAmount = 5

var (
	content = map[string]string{
		"first@gmaill.com": "12345Qwerty",
		"second@gmail.com": "hjtbgnjekmf",
		"third@gmail.com":  "eirugjlkwe",
		"fourth@mail.ru":   "erugjlkm",
		"fifth@yandex.ru":  "329789849nw",
	}
)

func createUsers() []*models.User {
	users := make([]*models.User, 0, usersAmount)

	for email, password := range content {
		user, err := models.NewUser(email, password)
		if err != nil {
			log.Fatalf("can't create new user: %s", err.Error())
		}

		users = append(users, user)
	}

	return users
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := createUsers()

	for _, user := range users {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Email, user.HashedPassword, user.Verificated).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.CreateUser(t.Context(), user)
		assert.Nil(t, err)
	}

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindUserByEmail(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := createUsers()

	for _, user := range users {
		row := mock.NewRows([]string{"user_id", "email", "hashed_password", "verificated"}).
			AddRow(user.ID, user.Email, user.HashedPassword, user.Verificated)

		mock.ExpectQuery("SELECT user_id, email, hashed_password, verificated FROM users").
			WithArgs(user.Email).WillReturnRows(row)

		u, err := repo.FindUserByEmail(t.Context(), user.Email)

		assert.Nil(t, err)
		assert.Equal(t, user.ID, u.ID)
		assert.Equal(t, user.Email, u.Email)
		assert.Equal(t, user.HashedPassword, u.HashedPassword)
		assert.Equal(t, user.Verificated, u.Verificated)
		assert.Nil(t, mock.ExpectationsWereMet())
	}

}

func TestVerifyEmail(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	users := createUsers()

	for _, user := range users {
		mock.ExpectExec("UPDATE users").
			WithArgs(user.Email).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.VerifyEmail(t.Context(), user.Email)

		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	}
}
