package app

import (
	"context"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/interfaces/http"
)

type App struct {
	server *http.Server
}

func (app *App) Start(ctx context.Context, healthCheckTime time.Duration) {
	app.server.Start(ctx, healthCheckTime)
}

func (app *App) Stop(ctx context.Context) {
	app.server.Close(ctx)
}
