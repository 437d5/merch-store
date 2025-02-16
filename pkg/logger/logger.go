package logger

import (
	"log/slog"
	"os"
)

func NewLogger(mode string, lvl slog.Level) *slog.Logger {
	var logger *slog.Logger

	switch mode {
	case "JSON":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		}))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		}))
	}

	slog.SetDefault(logger)

	slog.Info("Logger initialized")

	return logger
}
