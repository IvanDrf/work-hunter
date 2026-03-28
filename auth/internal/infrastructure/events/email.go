package events

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
)

type EmailWorker struct {
	consumer     service.EmailConsumer
	emailService service.EmailService
}

func NewEmailWorker(consumer service.EmailConsumer, emailService service.EmailService) *EmailWorker {
	return &EmailWorker{
		consumer:     consumer,
		emailService: emailService,
	}
}

func (w *EmailWorker) Start(ctx context.Context) {
	w.consumer.ProcessEmailsFromQueue(ctx, func(msg *models.EmailMessage) error {
		if !msg.IsTokenValid() {
			return models.Error{
				Message: "token is outdated",
				Code:    models.ErrOutdatedToken,
			}
		}

		return w.emailService.SendVerificationEmail(msg.Email, msg.Token)
	})
}

func (w *EmailWorker) Stop() {
	w.consumer.Close()
}
