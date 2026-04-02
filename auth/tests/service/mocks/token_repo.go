package mocks

import (
	"context"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
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
