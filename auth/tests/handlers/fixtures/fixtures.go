package fixtures

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

var (
	RegisterRequests = []*auth_api.User{
		{Email: "first@gmail.com", Password: "123456789", Role: string(models.EMPLOYEE)},
		{Email: "second@gmail.com", Password: "erjrglm", Role: string(models.EMPLOYEE)},
		{Email: "third@gmail.com", Password: "eroiigkml", Role: string(models.EMPLOYEE)},
		{Email: "fourth@gmail.com", Password: "eorigmke;r,", Role: string(models.EMPLOYEE)},
		{Email: "fifth@gmail.com", Password: "wlekfwef", Role: string(models.EMPLOYEE)},
	}

	UnregistredUsers = []*auth_api.User{
		{Email: "unreg@gmail.com", Password: "123456789", Role: string(models.EMPLOYEE)},
		{Email: "kjrngmr", Password: "erjrglm", Role: string(models.EMPLOYEE)},
		{Email: "unreg@main.ru", Password: "eroiigkml", Role: string(models.EMPLOYEE)},
		{Email: "reg4@mail.ru", Password: "eorigmke;r,", Role: string(models.EMPLOYEE)},
		{Email: "another@gmail.com", Password: "wlekfwef", Role: string(models.EMPLOYEE)},
	}

	InvalidRoleRequests = []*auth_api.User{
		{Email: "first@gmail.com", Password: "123456789", Role: "ekrjgnlmker,."},
		{Email: "second@gmail.com", Password: "erjrglm", Role: ""},
		{Email: "third@gmail.com", Password: "eroiigkml", Role: "erkger"},
		{Email: "fourth@gmail.com", Password: "eorigmke;r,", Role: "12332"},
		{Email: "fifth@gmail.com", Password: "wlekfwef", Role: "invalid_role"},
	}

	InvalidPasswordRequests = []*auth_api.User{
		{Email: "first@gmail.com", Password: "1234", Role: string(models.EMPLOYEE)},
		{Email: "second@gmail.com", Password: "kew2", Role: string(models.EMPLOYEE)},
		{Email: "third@gmail.com", Password: "eroiigkmleklrjghiuwjkefmwjefhuewiwdfmjasfnhuwf", Role: string(models.EMPLOYEE)},
		{Email: "fourth@gmail.com", Password: "", Role: string(models.EMPLOYEE)},
		{Email: "fifth@gmail.com", Password: "_w", Role: string(models.EMPLOYEE)},
	}

	InvalidEmailRequests = []*auth_api.User{
		{Email: "erjnglmke,r", Password: "123456789", Role: string(models.EMPLOYEE)},
		{Email: "", Password: "erjrglm", Role: string(models.EMPLOYEE)},
		{Email: "email.com", Password: "eroiigkml", Role: string(models.EMPLOYEE)},
		{Email: "12134", Password: "eorigmke;r,", Role: string(models.EMPLOYEE)},
		{Email: "printf", Password: "wlekfwef", Role: string(models.EMPLOYEE)},
	}
)
