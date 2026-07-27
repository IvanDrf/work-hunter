package clients

import (
	"context"
	"fmt"
	"log/slog"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userClient struct {
	retries int

	conn   *grpc.ClientConn
	client user_api.UserClient
}

func NewUserClient(host string, port int, retries int) *userClient {
	client, conn := connectToUser(host, port)
	return &userClient{
		retries: retries,
		conn:    conn,
		client:  client,
	}
}

func connectToUser(host string, port int) (user_api.UserClient, *grpc.ClientConn) {
	log := slog.With(slog.String("client", "user"))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("can't connect to user service", slog.String("error", err.Error()))
		return nil, nil
	}

	return user_api.NewUserClient(conn), conn
}

func (c *userClient) Close() {
	log := slog.With(slog.String("client", "user"))

	if c != nil {
		c.conn.Close()
	}

	log.Info("auth client is closed")
}

func (c *userClient) Health(ctx context.Context) {
	panic("UserClient not implemented Health")
}

func (c *userClient) SendGetCompanyName(ctx context.Context, userID string) (string, error) {
	panic("UserClient not implemented SendGetCompanyName")
}
