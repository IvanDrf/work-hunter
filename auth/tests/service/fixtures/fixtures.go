package fixtures

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/google/uuid"
)

// fixtures size, users, tokens
const size = 6

// fixtures
var (
	Users = map[string]string{
		"first@gmaill.com": "12345Qwerty",
		"second@gmail.com": "hjtbgnjekmf",
		"third@gmail.com":  "eirugjlkwe",
		"fourth@mail.ru":   "erugjlkm",
		"fifth@yandex.ru":  "329789849nw",
		"sixth@gmail.com":  "eorigjkplw",
	}

	// user ids
	UserIDs = [size]uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	UserIDsString = [size]string{
		UserIDs[0].String(),
		UserIDs[1].String(),
		UserIDs[2].String(),
		UserIDs[3].String(),
		UserIDs[4].String(),
		UserIDs[5].String(),
	}

	// unregistred users, using this in TestLoginUser
	Unregistered = map[string]string{
		"unregistred@gmail.com": "eruigjwkmelf",
		"un2egistred@gmail.com": "23iyrguhf",
		"unregthird@gmail.com":  "eorjgkml,",
		"unregfourth@mail.ru":   "wleojlfkm,",
		"unregfifth@yandex.ru":  "wojfkml;l,w.f",
		"unregsixth@gmail.com":  "wroefhjqk",
	}

	// Tokens for registred users
	Tokens = [size]string{
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
	}

	// invalid tokens
	UnusedTokens = [size]string{
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
		rules.GenerateToken(),
	}
)
