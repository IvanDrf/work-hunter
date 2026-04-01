package fixtures

import (
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

// fixtures
var (
	content = map[string]string{
		"first@gmaill.com": "12345Qwerty",
		"second@gmail.com": "hjtbgnjekmf",
		"third@gmail.com":  "eirugjlkwe",
		"fourth@mail.ru":   "erugjlkm",
		"fifth@yandex.ru":  "329789849nw",
	}
)

func CreateUsers() []*models.User {
	users := make([]*models.User, 0, len(content))

	for email, password := range content {
		user, err := models.NewUser(email, password)
		if err != nil {
			log.Fatalf("can't create new user: %s", err.Error())
		}

		users = append(users, user)
	}

	return users
}
