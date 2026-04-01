package models

import "github.com/google/uuid"

type JwtPayload struct {
	ID          uuid.UUID `json:"id"`
	Verificated bool      `json:"verificated"`
}
