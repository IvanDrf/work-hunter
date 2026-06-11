package models

import "github.com/google/uuid"

type JwtPayload struct {
	UserID      string `json:"user_id"`
	Verificated bool   `json:"verificated"`
	Role        Role   `json:"role"`
}

func (p *JwtPayload) IsPayloadValid() bool {
	_, err := uuid.Parse(p.UserID) // userID must be uuid
	_, ok := ROLES[string(p.Role)] // Role must be valid

	return err == nil && ok
}
