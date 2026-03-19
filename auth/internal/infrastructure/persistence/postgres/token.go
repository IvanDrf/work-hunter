package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (t *TokenRepo) CreateToken(ctx context.Context, email string, token *models.Token) error {
	const query = "INSERT INTO tokens (email, token, exp) VALUES($1, $2, $3) ON CONFLICT (email) DO UPDATE SET token = $3, exp = $3"

	_, err := t.db.ExecContext(ctx, query, email, token.Token, token.Exp, token.Token, token.Exp)
	return err
}

func (t *TokenRepo) FindEmailExpByToken(ctx context.Context, token string) (string, time.Time, error) {
	const query = "SELECT email, exp FROM tokens WHERE token = $1 LIMIT 1"

	email := ""
	exp := time.Time{}

	err := t.db.QueryRowContext(ctx, query, token).Scan(&email, &exp)
	return email, exp, err
}
