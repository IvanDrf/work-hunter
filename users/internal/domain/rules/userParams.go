package rules

import "errors"

// user status
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
	UserStatusDeleted  UserStatus = "deleted"
)

// user role
type UserRole string

const (
	UserRoleEmployee UserRole = "employee"
	UserRoleEmployer UserRole = "employer"
	UserRoleAdmin    UserRole = "admin"
)

func ValidateUserStatus(status UserStatus) error {
	if status != UserStatusActive && status != UserStatusInactive && status != UserStatusBlocked && status != UserStatusDeleted {
		return errors.New("invalid user status")
	}
	return nil
}

func ValidateUserRole(role UserRole) error {
	if role != UserRoleEmployee && role != UserRoleEmployer && role != UserRoleAdmin {
		return errors.New("invalid user role")
	}
	return nil
}
