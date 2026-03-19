package models

import (
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
)

const tokenTTL = 15 * time.Minute

type Token struct {
	Token string    `json:"token"`
	Exp   time.Time `json:"exp"`
}

func NewToken() *Token {
	return &Token{
		Token: rules.GenerateToken(),
		Exp:   time.Now().Add(tokenTTL).UTC(),
	}
}
