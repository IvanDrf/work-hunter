package messaging

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

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

func (c *RabbitMQConsumer) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *RabbitMQConsumer) ProcessEmailsFromQueue(ctx context.Context, fn func(*models.EmailMessage) error) {
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

			if err := processMessage(message, fn); err != nil {
				message.Reject(false)
				slog.Error("Error while processing message from broker", slog.String("error", err.Error()))
			} else {
				message.Ack(false)
			}
		}
	}
}

func processMessage(message rabbit.Delivery, fn func(*models.EmailMessage) error) error {
	email, err := parseMessage(message)
	if err != nil {
		return err
	}

	return fn(email)
}

func parseMessage(message rabbit.Delivery) (*models.EmailMessage, error) {
	email := &models.EmailMessage{}

	err := json.Unmarshal(message.Body, email)
	return email, err
}
