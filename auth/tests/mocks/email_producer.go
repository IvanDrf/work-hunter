package mocks

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type emailProducer struct {
	queue chan *models.EmailMessage
}

func NewEmailProducer(queue chan *models.EmailMessage) *emailProducer {
	return &emailProducer{
		queue: queue,
	}
}

func (e *emailProducer) SendEmailInQueue(ctx context.Context, message *models.EmailMessage) error {
	if e.queue == nil {
		return models.Error{
			Message: "queue is already closed",
			Code:    models.ErrCodeInternal,
		}
	}

	e.queue <- message
	return nil
}

func (e *emailProducer) Close() {
	close(e.queue)
	e.queue = nil
}
