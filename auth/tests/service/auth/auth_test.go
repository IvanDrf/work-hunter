package auth_test

import (
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/mocks"
)

func newAuthService() *service.AuthService {
	return service.NewAuthService(mocks.NewUserRepo(), mocks.Jwter)
}
