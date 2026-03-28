package rules

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const TokenTTL = 15 * time.Minute

func GenerateToken() string {
	buff := make([]byte, 32)

	rand.Read(buff)
	return hex.EncodeToString(buff)
}

func NewExpTime() time.Time {
	return time.Now().UTC().Add(TokenTTL)
}
