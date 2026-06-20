package config

import "fmt"

type RabbitMQConfig struct {
	RabbitMQHost string `env:"RABBITMQ_HOST"`
	RabbitMQPort int    `env:"RABBITMQ_PORT"`

	RabbitMQUsername string `env:"RABBITMQ_USER"`
	RabbitMQPassword string `env:"RABBITMQ_PASSWORD"`

	RabbitMQProducerQueue string `env:"RABBITMQ_PRODUCER_QUEUE"`
	RabbitMQConsumerQueue string `env:"RABBITMQ_CONSUMER_QUEUE"`
}

func (r *RabbitMQConfig) RABBITMQ_DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.RabbitMQUsername, r.RabbitMQPassword, r.RabbitMQHost, r.RabbitMQPort)
}
