package clients

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/IvanDrf/work-hunter/pkg/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	host string
	port int

	conn   *grpc.ClientConn
	client auth_api.AuthClient
}

func NewAuthClient(host string, port int) *authClient {
	c := &authClient{}
	c.connect()

	return c
}

func (c *authClient) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *authClient) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *authClient) connect() error {
	log := slog.With(slog.String("client", "auth"))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", c.host, c.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("can't connect to auth service", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't connect to auth servuce, auth is unavailable",
			Code:    models.ErrCodeInternal,
		}
	}

	c.client = auth_api.NewAuthClient(conn)
	c.conn = conn

	return nil
}

func (c *authClient) Close() {
	log := slog.With(slog.String("client", "auth"))

	if c != nil {
		c.conn.Close()
	}

	log.Info("auth client is closed")
}

func (c *authClient) SendRegisterRequest(ctx context.Context, email string, password string, role models.UserRole) (string, string, error) {
	log := slog.With(slog.String("client", "auth"))
	log.InfoContext(ctx, "register request", slog.String("email", email))

	if err := isConnected(c, c.connect); err != nil {
		slog.ErrorContext(ctx, "can't register user, auth service is unavailable", slog.String("email", email), slog.String("error", err.Error()))
		return "", "", err
	}

	resp, err := c.client.Register(ctx, &auth_api.User{
		Email:    email,
		Password: password,
		Role:     common.UserRole(role),
	})
	if err != nil {
		slog.ErrorContext(ctx, "can't register user, auth service returned error", slog.String("email", email), slog.String("error", err.Error()))
		return "", "", err
	}

	slog.InfoContext(ctx, "successfully registred user", slog.String("email", email))
	return resp.Access, resp.Refresh, nil
}
