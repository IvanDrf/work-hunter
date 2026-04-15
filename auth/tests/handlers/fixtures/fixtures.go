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
)

var (
	Users = []*auth_api.User{
		{Email: "first@gmail.com", Password: "123456789", Role: 0},
		{Email: "second@gmail.com", Password: "erjrglm", Role: 0},
		{Email: "third@gmail.com", Password: "eroiigkml", Role: 0},
		{Email: "fourth@gmail.com", Password: "eorigmke;r,", Role: 0},
		{Email: "fifth@gmail.com", Password: "wlekfwef", Role: 0},
	}

	UnregistredUsers = []*auth_api.User{
		{Email: "unreg@gmail.com", Password: "123456789", Role: 0},
		{Email: "kjrngmr", Password: "erjrglm", Role: 0},
		{Email: "unreg@main.ru", Password: "eroiigkml", Role: 0},
		{Email: "reg4@mail.ru", Password: "eorigmke;r,", Role: 0},
		{Email: "another@gmail.com", Password: "wlekfwef", Role: 0},
	}

	InvalidRoleRequests = []*auth_api.User{
		{Email: "first@gmail.com", Password: "123456789", Role: -1},
		{Email: "second@gmail.com", Password: "erjrglm", Role: -2},
		{Email: "third@gmail.com", Password: "eroiigkml", Role: 40},
		{Email: "fourth@gmail.com", Password: "eorigmke;r,", Role: 42},
		{Email: "fifth@gmail.com", Password: "wlekfwef", Role: 90},
	}

	InvalidPasswordRequests = []*auth_api.User{
		{Email: "first@gmail.com", Password: "1234", Role: 0},
		{Email: "second@gmail.com", Password: "kew2", Role: 0},
		{Email: "third@gmail.com", Password: "eroiigkmleklrjghiuwjkefmwjefhuewiwdfmjasfnhuwf", Role: 0},
		{Email: "fourth@gmail.com", Password: "", Role: 0},
		{Email: "fifth@gmail.com", Password: "_w", Role: 0},
	}

	InvalidEmailRequests = []*auth_api.User{
		{Email: "erjnglmke,r", Password: "123456789", Role: 0},
		{Email: "", Password: "erjrglm", Role: 0},
		{Email: "email.com", Password: "eroiigkml", Role: 0},
		{Email: "12134", Password: "eorigmke;r,", Role: 0},
		{Email: "printf", Password: "wlekfwef", Role: 0},
	}
)
