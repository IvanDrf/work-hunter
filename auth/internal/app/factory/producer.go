package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	messaging "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/messaging/rabbitmq"
)

func (f *Factory) newProducer() service.EmailProducer {
	return messaging.NewRabbitMqProducer(messaging.Connect(&f.cfg.Broker, f.cfg.Broker.ProducerQueue))
}
