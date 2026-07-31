package clients

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"github.com/IvanDrf/work-hunter/pkg/common"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	retries int

	conn   *grpc.ClientConn
	client auth_api.AuthClient
}

func NewAuthClient(host string, port int, retries int) *authClient {
	client, conn := connectToAuth(host, port)
	return &authClient{
		retries: retries,

		client: client,
		conn:   conn,
	}
}

func connectToAuth(host string, port int) (auth_api.AuthClient, *grpc.ClientConn) {
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
	ctx = adapters.InsertLogger(ctx, log)

	const healthRetries = 2

	resp, err := retry(ctx, healthRetries, func() (any, error) {
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
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
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
			Access:  resp.GetAccess(),
			Refresh: resp.GetRefresh(),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return resp.(*models.Tokens), nil
}

func (c *authClient) SendLoginRequest(ctx context.Context, email string, password string) (*models.Tokens, error) {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendLoginRequest"), slog.String("email", email))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
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
			Access:  resp.GetAccess(),
			Refresh: resp.GetRefresh(),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*models.Tokens), nil
}
func (c *authClient) SendChangePasswordRequest(ctx context.Context, access string, oldPassword string, newPassword string) error {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendChangePasswordRequest"))
	ctx = adapters.InsertLogger(ctx, log)

	_, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.ChangePassword(ctx, &auth_api.ChangePasswordRequest{
			Access: access,
			Old:    oldPassword,
			New:    newPassword,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't change password for user, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return resp, nil
	})

	return err
}

func (c *authClient) SendRefreshTokensRequest(ctx context.Context, refresh string) (*models.Tokens, error) {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendRefreshTokensRequest"))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.RefreshTokens(ctx, &auth_api.RefreshToken{
			Refresh: refresh,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't change password for user, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return &models.Tokens{
			Access:  resp.GetAccess(),
			Refresh: resp.GetRefresh(),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*models.Tokens), nil
}

func (c *authClient) SendIsTokenValidRequest(ctx context.Context, access string) (*models.TokenPayload, error) {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendIsTokenValidRequest"))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.IsTokenValid(ctx, &auth_api.AccessToken{
			Access: access,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't check is access token valid, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		id, err := uuid.Parse(resp.GetId())
		if err != nil {
			return nil, models.Error{
				Message: "invalid user_id in token, not uuid",
				Code:    models.ErrCodeInvalidArgument,
			}
		}

		return &models.TokenPayload{
			ID:          id,
			Verificated: resp.GetVerificated(),
			Role:        models.UserRole(resp.GetRole()),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*models.TokenPayload), nil
}

func (c *authClient) SendDeleteUserRequest(ctx context.Context, access string, password string) error {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "DeleteUser"))
	ctx = adapters.InsertLogger(ctx, log)

	_, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.DeleteUser(ctx, &auth_api.DeleteUserRequest{
			Access:   access,
			Password: password,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't delete user, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return resp, nil
	})

	return err
}

func (c *authClient) SendVerificationEmailRequest(ctx context.Context, access string) error {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendVerificationEmailRequest"))
	ctx = adapters.InsertLogger(ctx, log)

	_, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.SendVerificationEmail(ctx, &auth_api.AccessToken{
			Access: access,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't send verification email, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return resp, nil
	})

	return err
}

func (c *authClient) SendVerifyEmailRequest(ctx context.Context, token string) (*models.Tokens, error) {
	log := slog.With(slog.String("client", "auth"), slog.String("request", "SendVerifyEmailRequest"))
	ctx = adapters.InsertLogger(ctx, log)

	resp, err := retry(ctx, c.retries, func() (any, error) {
		resp, err := c.client.VerifyEmail(ctx, &auth_api.VerifToken{
			Token: token,
		})
		if err != nil {
			log.ErrorContext(ctx, "can't send verify user email, auth service returned error", slog.String("error", err.Error()))
			return nil, err
		}

		return &models.Tokens{
			Access:  resp.GetAccess(),
			Refresh: resp.GetRefresh(),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*models.Tokens), nil
}
