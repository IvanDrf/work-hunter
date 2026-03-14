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
	const query = "INSERT INTO users(user_id, email, hashed_password, verificated) VALUES($1, $2, $3)"

	_, err := a.db.ExecContext(ctx, query, user.ID, user.Email, user.HashedPassword, user.Verificated)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) FindUser(ctx context.Context, email string) (*models.User, error) {
	const query = "SELECT user_id, email, hashed_password, verificated FROM users WHERE username = $1 LIMIT 1"

	rows, err := a.db.QueryContext(ctx, query, email)
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
