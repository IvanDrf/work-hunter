package fixtures

import (
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

const (
	NewPassword        = "new password"
	InvalidPassword    = "invalid password"
	InvalidOldPassword = "invalid old password"
	InvalidNewPassword = ""

	InvalidUserID   = "invalid user id"
	InvalidUserRole = "invalid user role"

	InvalidVerificationToken = "jerngmkieemueqli-hefkw"
)

var (
	invalidRoles = [5]auth_api.Role{-1, -2, -3, -4, -50}
	emails       = [5]string{"first@gmail.com", "second@gmail.com", "third@gmail.com", "fourth@gmail.com", "fifth@gmail.com"}
	passwords    = [5]string{"123456789", "erjrglm", "eroiigkml", "eorigmke;r", "wlekfwef"}

	Users = []*auth_api.User{
		{Email: emails[0], Password: passwords[0], Role: 0},
		{Email: emails[1], Password: passwords[1], Role: 0},
		{Email: emails[2], Password: passwords[2], Role: 0},
		{Email: emails[3], Password: passwords[3], Role: 0},
		{Email: emails[4], Password: passwords[4], Role: 0},
	}

	UnregistredUsers = []*auth_api.User{
		{Email: "unreg@gmail.com", Password: passwords[0], Role: 0},
		{Email: "kjrngmr", Password: passwords[1], Role: 0},
		{Email: "unreg@main.ru", Password: passwords[2], Role: 0},
		{Email: "reg4@mail.ru", Password: passwords[3], Role: 0},
		{Email: "another@gmail.com", Password: passwords[4], Role: 0},
	}

	InvalidRoleRequests = []*auth_api.User{
		{Email: emails[0], Password: passwords[0], Role: invalidRoles[0]},
		{Email: emails[1], Password: passwords[1], Role: invalidRoles[1]},
		{Email: emails[2], Password: passwords[2], Role: invalidRoles[3]},
		{Email: emails[3], Password: passwords[3], Role: invalidRoles[3]},
		{Email: emails[4], Password: passwords[4], Role: invalidRoles[4]},
	}

	InvalidPasswordRequests = []*auth_api.User{
		{Email: emails[0], Password: "1234", Role: 0},
		{Email: emails[1], Password: "kew2", Role: 0},
		{Email: emails[2], Password: "eroiigkmleklrjghiuwjkefmwjefhuewiwdfmjasfnhuwf", Role: 0},
		{Email: emails[3], Password: "", Role: 0},
		{Email: emails[4], Password: "_w", Role: 0},
	}

	InvalidEmailRequests = []*auth_api.User{
		{Email: "erjnglmke,r", Password: passwords[0], Role: 0},
		{Email: "", Password: passwords[1], Role: 0},
		{Email: "email.com", Password: passwords[2], Role: 0},
		{Email: "12134", Password: passwords[3], Role: 0},
		{Email: "printf", Password: passwords[4], Role: 0},
	}
)
