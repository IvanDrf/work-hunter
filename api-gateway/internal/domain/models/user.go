package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	Role UserRole `json:"role"`
}

func (u *User) IsUserValid() bool {
	return u != nil && u.Email != "" && u.Password != ""
}
