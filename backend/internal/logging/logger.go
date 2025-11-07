package logging

import (
	"log/slog"
	"os"
)

func NewLogger(appName, env string) *slog.Logger {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler).With(
		"service", appName,
		"environment", env,
	)
}
