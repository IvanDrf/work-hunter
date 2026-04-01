package fixtures

import "github.com/IvanDrf/work-hunter/auth/internal/domain/rules"

// fixtures
var (
	Content = map[string]string{
		"second@gmail.com": rules.GenerateToken(),
		"first@gmaill.com": rules.GenerateToken(),
		"third@gmail.com":  rules.GenerateToken(),
		"fourth@mail.ru":   rules.GenerateToken(),
		"fifth@yandex.ru":  rules.GenerateToken(),
	}
)
