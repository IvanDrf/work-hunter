package rules

const (
	minPasswordLen = 5
	maxPasswordLen = 30
)

func IsPasswordCorrect(password string) bool {
	l := len(password)

	return minPasswordLen <= l && l <= maxPasswordLen
}
