package token_test

import (
	"testing"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	r "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/redis"
	"github.com/IvanDrf/work-hunter/auth/tests/repo/token/fixtures"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func connect() (*r.TokenRepo, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()

	repo := r.NewTokenRepo(client)
	return repo, mock
}

func TestCreateToken(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	for email, token := range fixtures.Content {
		mock.ExpectSet(token, email, rules.TokenTTL).RedisNil()
		err := repo.CreateToken(t.Context(), email, token, rules.TokenTTL)

		require.Equal(t, redis.Nil, err)
		require.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestFindEmailByToken(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	for email, token := range fixtures.Content {
		mock.ExpectGet(token).SetVal(email)

		e := repo.FindEmailByToken(t.Context(), token)

		assert.Equal(t, email, e)
		require.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestDeleteToken(t *testing.T) {
	t.Parallel()

	repo, mock := connect()
	defer repo.Close()

	for _, token := range fixtures.Content {
		mock.ExpectDel(token).RedisNil()

		err := repo.DeleteToken(t.Context(), token)

		assert.Equal(t, redis.Nil, err)
		require.NoError(t, mock.ExpectationsWereMet())
	}
}
