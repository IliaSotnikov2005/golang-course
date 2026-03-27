package sl

import (
	"log/slog"
	"os"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/config"
)

func SetupLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: getLogLevel(cfg),
		// AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func getLogLevel(cfg *config.Config) slog.Level {
	switch cfg.LogLevel {
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
