package postgres

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(cfg *config.DatabaseConfig) *AuthRepo {
	return &AuthRepo{
		db: Connect(cfg),
	}
}

func (a *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	const query = "INSERT INTO users(user_id, username, hashed_password) VALUES($1, $2, $3)"

	_, err := a.db.ExecContext(ctx, query, user.ID, user.Username, user.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) FindUser(ctx context.Context, username string) (*models.User, error) {
	const query = "SELECT user_id, username, hashed_password FROM users WHERE username = $1 LIMIT 1"

	rows, err := a.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := models.User{}
	err = rows.Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
