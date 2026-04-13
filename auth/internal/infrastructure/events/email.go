package events

import (
	"context"
	"log/slog"

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
	slog.Info("Starting WORKER")

	w.consumer.ProcessEmailsFromQueue(ctx, func(msg *models.EmailMessage) error {
		slog.Info("worker: got message from queue")
		if !msg.IsTokenValid() {
			slog.Warn("worker: token is outdated")
			return models.Error{
				Message: "token is outdated",
				Code:    models.ErrCodeOutdatedToken,
			}
		}

		return w.emailService.SendVerificationEmail(msg.Email, msg.Token)
	})
}

func (w *EmailWorker) Stop() {
	slog.Info("Stopping WORKER")
	w.consumer.Close()
}
