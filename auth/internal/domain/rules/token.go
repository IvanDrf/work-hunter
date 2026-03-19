package rules

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken() string {
	buff := make([]byte, 32)

	rand.Read(buff)
	return hex.EncodeToString(buff)
}
