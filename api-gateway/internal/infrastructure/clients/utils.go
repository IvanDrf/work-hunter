package clients

import (
	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"google.golang.org/grpc"
)

type Client interface {
	Conn() *grpc.ClientConn
	Address() string
}

func isConnected(client Client, retry func() error) error {
	if client.Conn() != nil {
		return nil
	}

	return retryToConnect(retry)
}

func retryToConnect(retry func() error) error {
	if err := retry(); err != nil {
		return models.Error{
			Message: "internal error, auth service is unavailable",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}
