package events

import "github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"

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

func (w *EmailWorker) Start() {

}
