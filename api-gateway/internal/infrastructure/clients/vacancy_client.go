package clients

import (
	"fmt"
	"log/slog"

	"context"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
	"github.com/IvanDrf/work-hunter/pkg/vacancy_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type vacancyClient struct {
	retries int

	conn   *grpc.ClientConn
	client vacancy_api.VacancyClient
}

func NewVacancyClient(host string, port int, retries int) *vacancyClient {
	client, conn := connectToVacancy(host, port)
	return &vacancyClient{
		retries: retries,
		conn:    conn,
		client:  client,
	}
}

func connectToVacancy(host string, port int) (vacancy_api.VacancyClient, *grpc.ClientConn) {
	log := slog.With(slog.String("client", "vacancy"))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("can't connect to vacancy service", slog.String("error", err.Error()))
		return nil, nil
	}

	return vacancy_api.NewVacancyClient(conn), conn
}

func (c *vacancyClient) Close() {
	log := slog.With(slog.String("client", "vacancy"))

	if c != nil {
		c.conn.Close()
	}

	log.Info("vacancy client is closed")
}

func (c *vacancyClient) Health(ctx context.Context) {
	log := slog.With(slog.String("client", "vacancy"))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.Health(ctx, nil)
		if err != nil {
			log.ErrorContext(ctx, "can't check vacancy service health, vacancy service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return resp, nil
	})

	if err != nil {
		log.ErrorContext(ctx, "vacancy service is not available now", slog.String("error", err.Error()))
	} else {
		log.InfoContext(ctx, "vacancy service is available now", slog.Any("resp", resp))
	}
}
