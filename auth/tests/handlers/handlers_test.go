package handlers_test

import (
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"

	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
)

func init() {
	slog.SetDefault(slog.New(slog.DiscardHandler))
}

func newHandlers(queue chan *models.EmailMessage) *handlers.Handler {
	userRepo := mocks.NewUserRepo()
	tokenRepo := mocks.NewTokenRepo()
	producer := mocks.NewEmailProducer(queue)

	return handlers.NewHandler(service.NewAuthService(userRepo, mocks.Jwter), service.NewVerificationService(producer, userRepo, tokenRepo, mocks.Jwter))
}
