package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	rabbit "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn  *rabbit.Connection
	ch    *rabbit.Channel
	queue *rabbit.Queue
}

func NewRabbitMqConsumer(cfg *config.RabbitMQConfig) *RabbitMQConsumer {
	conn, ch := connect(cfg)
	queue := declareQueue(cfg.ConsumerQueue, ch)

	return &RabbitMQConsumer{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}
}

func (c *RabbitMQConsumer) Stop() {
	c.ch.Close()
	c.conn.Close()
}

func (c *RabbitMQConsumer) GetEmailsFromQueue(ctx context.Context, output chan<- *models.EmailMessage) {
	messages, err := c.ch.Consume(c.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("can't start consuming messages from rabbitmq: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, ok := <-messages
			if !ok {
				continue
			}

			email, err := parseMessage(message)
			if err != nil {
				// add logging
				message.Reject(false)
				continue
			}

			output <- email
			message.Ack(false)
		}
	}
}

func parseMessage(message rabbit.Delivery) (*models.EmailMessage, error) {
	email := &models.EmailMessage{}

	err := json.Unmarshal(message.Body, email)
	return email, err
}
