package postgres

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) Close() {
	u.db.Close()
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	const query = "INSERT INTO users(user_id, email, hashed_password, verificated) VALUES($1, $2, $3)"

	_, err := u.db.ExecContext(ctx, query, user.ID, user.Email, user.HashedPassword, user.Verificated)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) FindUser(ctx context.Context, email string) (*models.User, error) {
	const query = "SELECT user_id, email, hashed_password, verificated FROM users WHERE username = $1 LIMIT 1"

	rows, err := u.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := models.User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.HashedPassword, &user.Verificated)

		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserRepo) VerifyEmail(ctx context.Context, email string) error {
	const query = "UPDATE users SET verificated = true WHERE email = $1"

	_, err := u.db.ExecContext(ctx, query, email)
	return err
}
