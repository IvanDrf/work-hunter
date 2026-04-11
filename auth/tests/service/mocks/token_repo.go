package mocks

import (
	"context"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
)

type tokenRepo struct {
	storage map[string]struct {
		email string
		ttl   time.Time
	}
}

func NewTokenRepo() *tokenRepo {
	return &tokenRepo{
		storage: map[string]struct {
			email string
			ttl   time.Time
		}{},
	}
}

// Returns new filled token repo with tokens for registres users from fixtures.Users
func NewFilledTokenRepo() *tokenRepo {
	tokenRepo := NewTokenRepo()
	i := 0
	for email := range fixtures.Users {
		tokenRepo.CreateToken(context.TODO(), email, fixtures.Tokens[i], rules.TokenTTL)
		i++
	}

	return tokenRepo
}

func (t *tokenRepo) CreateToken(ctx context.Context, email string, token string, ttl time.Duration) error {
	t.storage[token] = struct {
		email string
		ttl   time.Time
	}{
		email: email,
		ttl:   time.Now().UTC().Add(ttl),
	}

	return nil
}

func (t *tokenRepo) FindEmailByToken(ctx context.Context, token string) string {
	if content, ok := t.storage[token]; ok && time.Now().UTC().Before(content.ttl) {
		return content.email
	}

	return ""
}

func (t *tokenRepo) DeleteToken(ctx context.Context, token string) error {
	if _, ok := t.storage[token]; ok {
		t.storage[token] = struct {
			email string
			ttl   time.Time
		}{}
		return nil
	}

	return models.Error{
		Code:    models.ErrCodeInternal,
		Message: "can't delete token",
	}

}

func (t *tokenRepo) Close() {
	t.storage = nil
}
