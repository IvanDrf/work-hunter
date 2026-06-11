package models

type Role string

const (
	ADMIN    Role = "ADMIN"
	EMPLOYEE Role = "EMPLOYEE"
	EMPLOYER Role = "EMPLOYER"
)

var ROLES = map[string]Role{
	"ADMIN":    ADMIN,
	"EMPLOYEE": EMPLOYEE,
	"EMPLOYER": EMPLOYER,
}
