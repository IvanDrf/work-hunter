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
	slog.Info("worker: starting")

	w.consumer.ProcessEmailsFromQueue(ctx, func(msg *models.EmailMessage) error {
		slog.Info("worker: got message from queue")
		if !msg.IsTokenValid() {
			slog.Warn("worker: token is outdated")
			return models.Error{
				Message: "token is outdated",
				Code:    models.ErrCodeOutdatedToken,
			}
		}

		err := w.emailService.SendVerificationEmail(msg.Email, msg.Token)
		if err != nil {
			slog.Error("worker: can't send email", slog.String("error", err.Error()))
			return err
		}

		slog.Info("worker: sent email")
		return err
	})
}

func (w *EmailWorker) Stop() {
	slog.Info("worker: stopping")
	w.consumer.Close()
}
