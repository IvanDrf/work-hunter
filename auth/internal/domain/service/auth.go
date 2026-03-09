package service

type AuthService interface {
	RegisterUser(username string, password string) (string, string, error)
	LoginUser(username string, password string) (string, string, error)

	RefreshTokens(access string, refresh string) (string, string, error)
}
