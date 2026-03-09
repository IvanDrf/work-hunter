package postgres

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type AuthRepo struct {
	db *sql.DB
}

func (a *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	const query = "INSERT INTO users(user_id, username, hashed_password) VALUES($1, $2, $3)"

	_, err := a.db.ExecContext(ctx, query, user.ID, user.Username, user.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}
