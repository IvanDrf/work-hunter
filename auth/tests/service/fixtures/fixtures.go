package fixtures

import (
	"github.com/google/uuid"
)

// fixtures
var (
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
