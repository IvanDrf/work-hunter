package adapters

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func InitSlogLogger(level string, file *os.File) {
	opt := &slog.HandlerOptions{
		AddSource: isSourceNeeded(level),
		Level:     selectLevel(level),
	}

	switch {
	case strings.HasSuffix(file.Name(), "json"):
		slog.New(slog.NewJSONHandler(file, opt))
	case strings.HasSuffix(file.Name(), "txt") || strings.HasSuffix(file.Name(), "log"):
		slog.New(slog.NewTextHandler(file, opt))
	default:
		slog.New(slog.NewJSONHandler(os.Stdout, opt))
	}
}

func isSourceNeeded(level string) bool {
	return strings.ToLower(level) == "debug"
}

func selectLevel(level string) slog.Leveler {
	l := strings.ToLower(level)

	levels := map[string]slog.Leveler{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	if val, ok := levels[l]; ok {
		return val
	}

	log.Fatalf("invalid logger level, level=%s", level)
	return slog.LevelDebug
}
