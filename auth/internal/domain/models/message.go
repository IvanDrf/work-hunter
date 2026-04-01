package models

import "time"

type EmailMessage struct {
	Email string `json:"email"`

	Token string    `json:"token"`
	Exp   time.Time `json:"exp"`
}

func (e *EmailMessage) IsTokenValid() bool {
	return time.Now().UTC().Before(e.Exp)
}
