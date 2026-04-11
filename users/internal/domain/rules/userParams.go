package rules

import "errors"

// user status
type UserStatus string

const (
	UserStatusActive   UserStatus = "USER_STATUS_ACTIVE"
	UserStatusInactive UserStatus = "USER_STATUS_INACTIVE"
	UserStatusBlocked  UserStatus = "USER_STATUS_BLOCKED"
	UserStatusDeleted  UserStatus = "USER_STATUS_DELETED"
)

// user role
type UserRole string

const (
	UserRoleUser      UserRole = "USER_ROLE_USER"
	UserRoleModerator UserRole = "USER_ROLE_MODERATOR"
	UserRoleAdmin     UserRole = "USER_ROLE_ADMIN"
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
