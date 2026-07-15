package app

import (
	"context"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/interfaces/http"
)

type App struct {
	server *http.Server
}

func (app *App) Start(ctx context.Context) {
	app.server.Start(ctx)
}

func (app *App) Stop(ctx context.Context) {
	app.server.Close(ctx)
}
