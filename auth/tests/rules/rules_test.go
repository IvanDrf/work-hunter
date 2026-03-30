package rules_test

import (
	"sync"
	"testing"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/auth/tests/rules/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	t.Parallel()

	wg := new(sync.WaitGroup)
	wg.Go(func() {
		for _, email := range fixtures.ValidEmails {
			assert.True(t, rules.IsEmailValid(email))
		}
	})

	wg.Go(func() {
		for _, email := range fixtures.InvalidEmails {
			assert.False(t, rules.IsEmailValid(email))
		}
	})

	wg.Wait()
}

func TestPasswordHashing(t *testing.T) {
	t.Parallel()

	for _, password := range fixtures.ValidPasswords {
		hashed, err := rules.HashPassword(password)

		assert.Nil(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
		assert.True(t, rules.IsPasswordsSame(password, hashed))
	}
}

func TestIsPasswordCorrect(t *testing.T) {
	t.Parallel()

	wg := new(sync.WaitGroup)
	wg.Go(func() {
		for _, password := range fixtures.ValidPasswords {
			assert.True(t, rules.IsPasswordCorrect(password))
		}
	})

	wg.Go(func() {
		for _, password := range fixtures.InvalidPasswords {
			assert.False(t, rules.IsPasswordCorrect(password))
		}
	})

	wg.Wait()
}

func TestGenerateToken(t *testing.T) {
	t.Parallel()

	const tokensAmount = 5
	// all tokens should be unique
	tokens := make(map[string]struct{}, tokensAmount)

	for range tokensAmount {
		token := rules.GenerateToken()
		assert.NotEmpty(t, token)

		_, ok := tokens[token]
		assert.False(t, ok)

		tokens[token] = struct{}{}
	}
}

func TestNewExpTimeForToken(t *testing.T) {
	t.Parallel()

	exp := time.Now().UTC().Add(rules.TokenTTL)
	actualExp := rules.NewExpTime()

	assert.Equal(t, exp.Year(), actualExp.Year())
	assert.Equal(t, exp.Day(), actualExp.Day())
	assert.Equal(t, exp.Minute(), actualExp.Minute())
}
