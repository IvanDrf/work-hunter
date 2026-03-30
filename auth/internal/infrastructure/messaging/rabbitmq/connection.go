package messaging

import (
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	rabbit "github.com/rabbitmq/amqp091-go"
)

func Connect(cfg *config.RabbitMQConfig, queueName string) (*rabbit.Connection, *rabbit.Channel, *rabbit.Queue) {
	conn, err := rabbit.Dial(cfg.DSN())
	if err != nil {
		log.Fatalf("can't connect to rabbitmq: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("can't open new channel in rabbitmq: %s", err)
	}

	queue := declareQueue(queueName, ch)

	return conn, ch, queue
}

func declareQueue(name string, ch *rabbit.Channel) *rabbit.Queue {
	queue, err := ch.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("can't declare queue in rabbitmq: %s", err)
	}

	return &queue
}
