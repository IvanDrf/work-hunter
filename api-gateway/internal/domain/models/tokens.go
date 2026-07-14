package models

import "github.com/google/uuid"

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type TokenPayload struct {
	ID          uuid.UUID `json:"id"`
	Verificated bool      `json:"verificated"`
	Role        UserRole  `json:"role"`
}
