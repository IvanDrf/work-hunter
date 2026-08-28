package config

type RabbitMQConfig struct {
	RabbitMQHost          string `env:"RABBITMQ_HOST"`
	RabbitMQPort          int    `env:"RABBITMQ_PORT"`
	RabbitMQUsername      string `env:"RABBITMQ_USER"`
	RabbitMQPassword      string `env:"RABBITMQ_PASSWORD"`
	RabbitMQConsumerQueue string `env:"RABBITMQ_CONSUMER_QUEUE"`
}
