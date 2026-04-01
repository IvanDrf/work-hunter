package service

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type EmailService interface {
	SendVerificationEmail(email string, token string) error
}

type EmailProducer interface {
	SendEmailInQueue(ctx context.Context, message *models.EmailMessage) error

	Close()
}

type EmailConsumer interface {
	ProcessEmailsFromQueue(ctx context.Context, fn func(msg *models.EmailMessage) error)

	Close()
}
