package rules

import "regexp"

func IsEmailValid(email string) bool {
	const re = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	return regexp.MustCompile(re).MatchString(email)
}
