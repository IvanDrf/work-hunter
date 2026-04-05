package verification_test

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/fixtures"
	"github.com/IvanDrf/work-hunter/auth/tests/service/mocks"
)

// email producer queue size
const size = 100

// email producer queue
var Queue = make(chan *models.EmailMessage, size)

func newVerificationService() *service.VerificationService {
	// add some users in mock database
	userRepo := mocks.NewUserRepo()
	i := 0
	for email, password := range fixtures.Users {
		userRepo.CreateUser(context.TODO(), &models.User{
			ID:             fixtures.UserIDs[i],
			Email:          email,
			HashedPassword: password,
			Verificated:    false,
		})
	}

	return service.NewVerificationService(mocks.NewEmailProducer(Queue), userRepo, mocks.NewTokenRepo(), mocks.Jwter)
}
