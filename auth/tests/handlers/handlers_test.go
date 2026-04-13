package handlers_test

import (
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"

	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
)

func init() {
	slog.SetDefault(slog.New(slog.DiscardHandler))
}

func newHandlers() *handlers.Handler {
	return handlers.NewHandler(service.NewAuthService(mocks.NewUserRepo(), mocks.Jwter), mocks.NewVerificationService())
}
