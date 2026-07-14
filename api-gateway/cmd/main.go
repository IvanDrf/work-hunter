package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/config"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/clients"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/interfaces/http"
)

func main() {
	config := config.LoadFromENV()
	auth := clients.NewAuthClient(config.Auth.Host, config.Auth.Port, config.App.Retries)

	s := http.NewServer(config.App.Host, config.App.Port, http.NewHandlers(auth, config.App.RequestTime), http.NewAuthMiddleware(auth, config.App.RequestTime))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()

	c, can := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer can()

	s.Close(c)
}
