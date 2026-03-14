package messaging

import (
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	rabbit "github.com/rabbitmq/amqp091-go"
)

func connect(cfg *config.RabbitMQConfig) (*rabbit.Connection, *rabbit.Channel) {
	conn, err := rabbit.Dial(cfg.DSN())
	if err != nil {
		log.Fatalf("can't connect to rabbitmq: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("can't open new channel in rabbitmq: %s", err)
	}

	return conn, ch
}
