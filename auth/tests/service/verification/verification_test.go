package verification_test

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/auth/tests/service/mocks"
)

// email producer queue size
const size = 100

// email producer queue
var Queue chan *models.EmailMessage = nil

func newVerificationService() *service.VerificationService {
	if Queue == nil {
		Queue = make(chan *models.EmailMessage, size)
	}

	userRepo, tokenRepo := mocks.NewFilledUserRepo(), mocks.NewFilledTokenRepo()
	return service.NewVerificationService(mocks.NewEmailProducer(Queue), userRepo, tokenRepo, mocks.Jwter)
}
