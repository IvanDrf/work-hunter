package fixtures

import (
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/auth/pkg"
	"github.com/google/uuid"
)

// fixtures
var (
	// Secret - valid Secret for  auth jwt tokens, auth service use this Secret
	Secret = rules.GenerateToken()

	// InvalidSecret - invalid secret for auth jwt tokens, auth service doesn't use this secret
	InvalidSecret = rules.GenerateToken()

	Jwter        = pkg.NewJwt(Secret, 1*time.Minute, 2*time.Minute)
	InvalidJwter = pkg.NewJwt(InvalidSecret, 1*time.Minute, 2*time.Minute)

	Users = map[string]string{
		"first@gmaill.com": "12345Qwerty",
		"second@gmail.com": "hjtbgnjekmf",
		"third@gmail.com":  "eirugjlkwe",
		"fourth@mail.ru":   "erugjlkm",
		"fifth@yandex.ru":  "329789849nw",
	}

	// user ids
	UserIDs = []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	// unregistred users, using this in TestLoginUser
	Unregistered = map[string]string{
		"unregistred@gmail.com": "eruigjwkmelf",
		"un2egistred@gmail.com": "23iyrguhf",
	}
)
