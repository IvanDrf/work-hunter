package config

import "fmt"

type RabbitMQConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`

	ProducerQueue string `yaml:"producer_queue"`
	ConsumerQueue string `yaml:"consumer_queue"`
}

func (r *RabbitMQConfig) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.Username, r.Password, r.Host, r.Port)
}
