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

var UserStatusInDB = map[string]string{"USER_STATUS_ACTIVE": "active", "USER_STATUS_INACTIVE": "inactive", "USER_STATUS_BLOCKED": "blocked"}

// user role
type UserRole string

const (
	UserRoleEmployee UserRole = "USER_ROLE_EMPLOYEE"
	UserRoleEmployer UserRole = "USER_ROLE_EMPLOYER"
	UserRoleAdmin    UserRole = "USER_ROLE_ADMIN"
)

var UserRoleInDB = map[string]string{"USER_ROLE_EMPLOYEE": "employee", "USER_ROLE_EMPLOYER": "employer", "USER_ROLE_ADMIN": "admin"}

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
