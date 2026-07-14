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
	client, conn := connect(host, port)
	return &authClient{
		host:    host,
		port:    port,
		retries: retries,

		client: client,
		conn:   conn,
	}
}

func connect(host string, port int) (auth_api.AuthClient, *grpc.ClientConn) {
	log := slog.With(slog.String("client", "auth"))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("can't connect to auth service", slog.String("error", err.Error()))
		return nil, nil
	}

	return auth_api.NewAuthClient(conn), conn
}

func (c *authClient) Close() {
	log := slog.With(slog.String("client", "auth"))

	if c != nil {
		c.conn.Close()
	}

	log.Info("auth client is closed")
}

func (c *authClient) Health(ctx context.Context) {
	log := slog.With(slog.String("client", "auth"))

	resp, err := retry(ctx, c.retries, log, func() (any, error) {
		resp, err := c.client.Health(ctx, nil)
		if err != nil {
			log.ErrorContext(ctx, "can't check auth service health, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return resp, nil
	})

	if err != nil {
		log.ErrorContext(ctx, "auth service is not available now", slog.String("error", err.Error()))
	} else {
		log.InfoContext(ctx, "auth service is available now", slog.Any("resp", resp))
	}
}

func (c *authClient) SendRegisterRequest(ctx context.Context, email string, password string, role models.UserRole) (*models.Tokens, error) {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendRegisterRequest"), slog.String("email", email))

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
		log.InfoContext(ctx, "successfully registred user")
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

func (c *authClient) SendLoginRequest(ctx context.Context, email string, password string) (*models.Tokens, error) {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendLoginRequest"), slog.String("email", email))

	resp, err := retry(ctx, c.retries, log, func() (any, error) {
		resp, err := c.client.Login(ctx, &auth_api.User{
			Email:    email,
			Password: password,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't login user, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		log.InfoContext(ctx, "successfully login user")
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
