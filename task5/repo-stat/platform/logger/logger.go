package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	LogLevel string `yaml:"log_level" default:"DEBUG"`
}

func MakeLogger(logLevel string) (*slog.Logger, error) {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("unknown log level: %s", logLevel)
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler), nil
}
