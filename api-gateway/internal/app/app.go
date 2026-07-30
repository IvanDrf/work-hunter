package app

import (
	"context"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/interfaces/http"
)

type App struct {
	server *http.Server
}

func (app *App) Start(ctx context.Context, periodCheckHealthTime time.Duration) {
	app.server.Start(ctx, periodCheckHealthTime)
}

func (app *App) Stop(ctx context.Context) {
	app.server.Close(ctx)
}
