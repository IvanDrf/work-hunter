package messaging

import rabbit "github.com/rabbitmq/amqp091-go"

type RabbitMQConsumer struct {
	conn  *rabbit.Connection
	ch    *rabbit.Channel
	queue *rabbit.Queue
}

func NewRabbitMqConsumer() {

}
