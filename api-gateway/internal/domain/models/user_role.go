package models

type UserRole int32

const (
	UNSPECIFIED UserRole = 0
	ADMIN       UserRole = 1
	EMPLOYEE    UserRole = 2
	EMPLOYER    UserRole = 3
)
