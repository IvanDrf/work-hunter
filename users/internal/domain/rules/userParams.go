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
	UserRoleUser      UserRole = "user"
	UserRoleModerator UserRole = "moderator"
	UserRoleAdmin     UserRole = "admin"
)

func ValidateUserStatus(status UserStatus) error {
	if status != UserStatusActive && status != UserStatusInactive && status != UserStatusBlocked && status != UserStatusDeleted {
		return errors.New("invalid user status")
	}
	return nil
}

func ValidateUserRole(role UserRole) error {
	if role != UserRoleUser && role != UserRoleModerator && role != UserRoleAdmin {
		return errors.New("invalid user role")
	}
	return nil
}
