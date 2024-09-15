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

	correlationID := ctx.Value(CtxKeyCorrelationID)
	if correlationID != nil {
		logger = logger.With(
			slog.String(string(CtxKeyCorrelationID), correlationID.(string)),
		)
	}

	clientID := ctx.Value(CtxKeyClientID)
	if clientID != nil {
		logger = logger.With(
			slog.String(string(CtxKeyClientID), clientID.(string)),
		)
	}

	return logger
}
