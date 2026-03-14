package logger

import (
	"log/slog"
	"os"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/config"
)

type Logger struct {
	*slog.Logger
}

func New(cfg *config.LoggerConfig) *Logger {
	var handler slog.Handler

	options := slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &options)
	default:
		handler = slog.NewTextHandler(os.Stdout, &options)
	}

	logger := slog.New(handler)

	return &Logger{logger}
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
