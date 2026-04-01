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
	const query = "INSERT INTO users(user_id, email, hashed_password, verificated) VALUES($1, $2, $3, $4)"

	_, err := u.db.ExecContext(ctx, query, user.ID, user.Email, user.HashedPassword, user.Verificated)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = "SELECT user_id, email, hashed_password, verificated FROM users WHERE email = $1 LIMIT 1"

	user := models.User{}
	err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.Verificated)

	return &user, err
}

func (u *UserRepo) VerifyEmail(ctx context.Context, email string) error {
	const query = "UPDATE users SET verificated = true WHERE email = $1"

	_, err := u.db.ExecContext(ctx, query, email)
	return err
}
