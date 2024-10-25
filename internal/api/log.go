package api

import (
	"context"
	"log/slog"
	"os"
)

var baseLogger = slog.New(
	slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	),
)

func Logger(ctx context.Context) *slog.Logger {
	logger := baseLogger.With()

	meta := NewRequestMeta(ctx)

	if meta.CorrelationID != "" {
		logger = logger.With(
			slog.String(string(ctxKeyCorrelationID), meta.CorrelationID),
		)
	}

	if meta.ClientID != "" {
		logger = logger.With(
			slog.String(string(ctxKeyClientID), meta.ClientID),
		)
	}

	return logger
}
