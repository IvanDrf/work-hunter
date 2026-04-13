package auth_test

import (
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/mocks"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelError)
}

func newAuthService() *service.AuthService {
	return service.NewAuthService(mocks.NewUserRepo(), mocks.Jwter)
}
