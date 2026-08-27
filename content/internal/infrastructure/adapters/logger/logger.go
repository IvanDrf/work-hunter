package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case "local", "dev":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return slog.New(handler)
}
