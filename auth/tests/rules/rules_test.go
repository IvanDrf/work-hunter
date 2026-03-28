package rules_test

import (
	"sync"
	"testing"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	t.Parallel()

	valids := []string{
		"test@gmail.com",
		"y@yandex.ru",
		"user@mail.ru",
	}

	invalids := []string{
		"41235612873",
		"user.com",
		"invalid@",
	}

	wg := new(sync.WaitGroup)
	wg.Go(func() {
		for _, email := range valids {
			assert.True(t, rules.IsEmailValid(email))
		}
	})

	wg.Go(func() {
		for _, email := range invalids {
			assert.False(t, rules.IsEmailValid(email))
		}
	})

	wg.Wait()
}

func TestPasswordHashing(t *testing.T) {
	t.Parallel()

	passwords := []string{
		"qwerty123",
		"printsf.f_mtA",
		"ersmkruwbrlnh123",
	}

	for _, password := range passwords {
		hashed, err := rules.HashPassword(password)

		assert.Nil(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
		assert.True(t, rules.IsPasswordsSame(password, hashed))
	}
}

func TestIsPasswordCorrect(t *testing.T) {
	t.Parallel()

	valids := []string{
		"qwerty123",
		"strongPASsWorD",
		"nu2msq",
	}

	invalids := []string{
		"rthjekpofwpoifjliuwekfmwjefwuikfmwel;f,wiejri",
		"1",
		"q",
		"wdf",
	}

	wg := new(sync.WaitGroup)
	wg.Go(func() {
		for _, password := range valids {
			assert.True(t, rules.IsPasswordCorrect(password))
		}
	})

	wg.Go(func() {
		for _, password := range invalids {
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
