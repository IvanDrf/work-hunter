package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	messaging "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/messaging/rabbitmq"
)

func (f *Factory) newConsumer() service.EmailConsumer {
	return messaging.NewRabbitMqConsumer(messaging.Connect(&f.cfg.Broker, f.cfg.Broker.ConsumerQueue))
}
