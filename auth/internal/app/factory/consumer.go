package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	messaging "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/messaging/rabbitmq"
)

func (f *Factory) NewConsumer() service.EmailConsumer {
	return messaging.NewRabbitMqConsumer(&f.cfg.Broker)
}
