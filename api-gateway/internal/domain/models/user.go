package models

type UserRole int32

const (
	UNSPECIFIED UserRole = 0
	ADMIN       UserRole = 1
	EMPLOYEE    UserRole = 2
	EMPLOYER    UserRole = 3
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	Role UserRole `json:"role"`
}

func (u *User) IsValid() bool {
	return u != nil && u.Email != "" && u.Password != ""
}

type UserInfo struct {
	Role        UserRole `json:"user_role"`
	UserID      string   `json:"user_id"`
	Verificated bool     `json:"verificated"`
}
