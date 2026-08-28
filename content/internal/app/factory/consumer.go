package factory

import (
	"fmt"

	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/messaging"
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/content/internal/infrastructure/clients/rabbitmq"
)

func (f *Factory) NewConsumer() (messaging.EventConsumer, error) {
	repo, err := f.newRepos()
	if err != nil {
		return nil, err
	}

	service := f.newContentService(repo)

	return f.newRabbitMQConsumer(service)
}

func (f *Factory) newRabbitMQConsumer(service service.ContentService) (*rabbitmq.UserConsumer, error) {
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		f.cfg.RabbitMQUsername,
		f.cfg.RabbitMQPassword,
		f.cfg.RabbitMQHost,
		f.cfg.RabbitMQPort,
	)

	return rabbitmq.NewUserConsumer(rabbitURL, f.cfg.RabbitMQConsumerQueue, service, f.log)
}
