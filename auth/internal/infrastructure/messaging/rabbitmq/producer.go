package messaging

import (
	"context"
	"encoding/json"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	rabbit "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn  *rabbit.Connection
	ch    *rabbit.Channel
	queue *rabbit.Queue
}

func NewRabbitMqProducer(conn *rabbit.Connection, ch *rabbit.Channel, queue *rabbit.Queue) *RabbitMQProducer {
	return &RabbitMQProducer{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}
}

func (p *RabbitMQProducer) Close() {
	p.ch.Close()
	p.conn.Close()
}

func (p *RabbitMQProducer) SendEmailInQueue(ctx context.Context, msg *models.EmailMessage) error {
	body, err := json.Marshal(*msg)
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx, "", p.queue.Name, false, false, rabbit.Publishing{
		Body: body,
	})
	if err != nil {
		return err
	}

	return nil
}
