package logging

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	var handler slog.Handler

	// Set default log level to info
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Log json in prod and pretty text elsewhere
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// In non-prod, set log level to debug
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
