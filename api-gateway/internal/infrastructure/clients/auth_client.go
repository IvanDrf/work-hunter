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
	host    string
	port    int
	retries int

	conn   *grpc.ClientConn
	client auth_api.AuthClient
}

func NewAuthClient(host string, port int, retries int) *authClient {
	c := &authClient{
		host:    host,
		port:    port,
		retries: retries,
	}
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

	conn, err := grpc.NewClient(c.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (c *authClient) SendRegisterRequest(ctx context.Context, email string, password string, role models.UserRole) (*models.Tokens, error) {
	log := slog.With(slog.String("client", "auth"))
	log.InfoContext(ctx, "register request", slog.String("email", email))

	resp, err := retry(ctx, c.retries, log, func() (any, error) {
		resp, err := c.client.Register(ctx, &auth_api.User{
			Email:    email,
			Password: password,
			Role:     common.UserRole(role),
		})
		if err != nil {
			log.ErrorContext(ctx, "can't register new user, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}
		slog.InfoContext(ctx, "successfully registred user", slog.String("email", email))
		return &models.Tokens{
			Access:  resp.Access,
			Refresh: resp.Refresh,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return resp.(*models.Tokens), nil
}
