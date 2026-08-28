package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IvanDrf/work-hunter/content/internal/domain/events"
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	service   service.ContentService
	log       *slog.Logger
}

func NewUserConsumer(url, queueName string, svc service.ContentService, log *slog.Logger) (*UserConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := ch.ExchangeDeclare("user_events", "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(q.Name, "user.deleted", "user_events", false, nil)
	if err != nil {
		return nil, err
	}

	return &UserConsumer{conn: conn, channel: ch, queueName: q.Name, service: svc, log: log}, nil
}

func (c *UserConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(c.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Info("stopping rabbitmq consumer")
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				var evt events.UserDeletedEvent
				if err := json.Unmarshal(d.Body, &evt); err != nil {
					c.log.Error("failed to unmarshal UserDeletedEvent", slog.String("error", err.Error()))
					continue
				}

				c.log.Info("received UserDeletedEvent", slog.String("userID", evt.UserID))
				_ = c.service.DeleteAllUserContent(ctx, evt.UserID)
			}
		}
	}()

	return nil
}

func (c *UserConsumer) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
