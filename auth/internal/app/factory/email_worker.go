package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/events"
	e "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/events"
)

func (f *Factory) NewEmailWorker() events.EmailWorker {
	consumer := f.newConsumer()
	smtp := f.newSmtpEmailService()

	return e.NewEmailWorker(consumer, smtp)
}
