package messaging

import (
	"log"

	rabbit "github.com/rabbitmq/amqp091-go"
)

func declareQueue(name string, ch *rabbit.Channel) *rabbit.Queue {
	queue, err := ch.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("can't declare queue in rabbitmq: %s", err)
	}

	return &queue
}
