package handlers_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"

	"github.com/IvanDrf/work-hunter/auth/tests/handlers/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/mocks"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.SetLogLoggerLevel(slog.LevelError)
}

func newHandlers(queue chan *models.EmailMessage) *handlers.Handler {
	userRepo := mocks.NewUserRepo()
	tokenRepo := mocks.NewTokenRepo()
	producer := mocks.NewEmailProducer(queue)

	return handlers.NewHandler(service.NewAuthService(userRepo, mocks.Jwter), service.NewVerificationService(producer, userRepo, tokenRepo, mocks.Jwter))
}

type Tokens struct {
	Access  string
	Refresh string
}

// Register users with handlers
func registerUsers(handlers *handlers.Handler, requests []*auth_api.User) map[string]*Tokens {
	tokens := make(map[string]*Tokens, len(requests))

	for _, req := range requests {
		resp, _ := handlers.Register(context.TODO(), req)

		tokens[req.Email] = &Tokens{
			Access:  resp.Access,
			Refresh: resp.Refresh,
		}
	}

	return tokens
}

func newHandlersWithVerifUsers(queue chan *models.EmailMessage) *handlers.Handler {
	users := make(map[string]string, len(fixtures.Users))
	for _, user := range fixtures.Users {
		users[user.Email] = user.Password
	}

	userRepo := mocks.NewFilledUserRepo(true, users)
	tokenRepo := mocks.NewTokenRepo()
	producer := mocks.NewEmailProducer(queue)

	return handlers.NewHandler(service.NewAuthService(userRepo, mocks.Jwter), service.NewVerificationService(producer, userRepo, tokenRepo, mocks.Jwter))
}
