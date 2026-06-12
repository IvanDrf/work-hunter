package models

import (
	"time"

	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Email       string    `db:"email"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name" json:"last_name"`
	CompanyName string    `db:"company_name"`

	Status rules.UserStatus `db:"status"`
	Role   rules.UserRole   `db:"role"`

	Verificated bool `db:"verificated"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(id uuid.UUID, email string, firstName, lastName, companyName *string, verificated bool) *User {
	u := &User{
		ID:    id,
		Email: email,

		Status: rules.UserStatusActive,
		Role:   rules.UserRoleUser,

		Verificated: verificated,
	}

	if firstName != nil {
		u.FirstName = *firstName
	}
	if lastName != nil {
		u.LastName = *lastName
	}
	if companyName != nil {
		u.CompanyName = *companyName
	}

	return u
}

func (u *User) UpdateUser(firstName, lastName, companyName *string, verificated *bool) {
	if firstName != nil {
		u.FirstName = *firstName
	}
	if lastName != nil {
		u.LastName = *lastName
	}
	if companyName != nil {
		u.CompanyName = *companyName
	}

	if verificated != nil {
		u.Verificated = *verificated
	}

}
