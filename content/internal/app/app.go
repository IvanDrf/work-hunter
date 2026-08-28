package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/IvanDrf/work-hunter/content/internal/app/factory"
	"github.com/IvanDrf/work-hunter/content/internal/config"
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/messaging"
	"github.com/IvanDrf/work-hunter/content/internal/interfaces/http/server/middleware"
	"github.com/gin-gonic/gin"
)

type App struct {
	cfg        *config.Config
	log        *slog.Logger
	httpServer *http.Server
	consumer   messaging.EventConsumer
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	factory := factory.NewFactory(cfg, log)

	consumer, err := factory.NewConsumer()
	if err != nil {
		log.Warn("rabbitmq connection failed, working without consumer", slog.String("error", err.Error()))
	}

	contentHandler, err := factory.NewHandler()
	if err != nil {
		return nil, err
	}

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	api.Use(middleware.RequireAPIKey(cfg.APIKey))

	api.POST("/content/resume/:user_id", contentHandler.UploadResume)
	api.POST("/content/avatar/:user_id", contentHandler.UploadAvatar)
	api.GET("/content/:type/:user_id", contentHandler.Download)
	api.DELETE("/content/:type/:user_id", contentHandler.Delete)

	httpServer := &http.Server{
		Addr:    cfg.Host + ":" + fmt.Sprint(cfg.Port),
		Handler: router,
	}

	return &App{cfg: cfg, log: log, httpServer: httpServer, consumer: consumer}, nil
}

func (a *App) Run(ctx context.Context) error {
	if a.consumer != nil {
		if err := a.consumer.Start(ctx); err != nil {
			a.log.Error("failed to start rabbitmq consumer", slog.String("error", err.Error()))
		}
	}

	go func() {
		a.log.Info("http server starting", slog.Int("port", a.cfg.Port))
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("http server failed", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()
	a.log.Info("shutting down application gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.log.Error("http server forced shutdown", slog.String("error", err.Error()))
	}
	if a.consumer != nil {
		a.consumer.Close()
	}

	a.log.Info("application stopped gracefully")
	return nil
}
