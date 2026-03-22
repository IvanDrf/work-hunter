package adapters

import (
	"log"
	"log/slog"
	"os"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
)

type LoggerLevel = string

const (
	debugLevel LoggerLevel = "debug"
	infoLevel  LoggerLevel = "info"
	warnLevel  LoggerLevel = "warn"
	errorLevel LoggerLevel = "error"
)

func InitLogger(cfg *config.AppConfig) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: selectSource(cfg.LoggerLevel),
		Level:     selectLevel(cfg.LoggerLevel),
	}))

	slog.SetDefault(logger)
}

func selectSource(level string) bool {
	return level == debugLevel
}

func selectLevel(level string) slog.Leveler {
	switch level {
	case debugLevel:
		return slog.LevelDebug

	case infoLevel:
		return slog.LevelInfo

	case warnLevel:
		return slog.LevelWarn

	case errorLevel:
		return slog.LevelError

	default:
		log.Fatalf("invalid logger level in cfg: %s", level)
	}

	return nil
}
