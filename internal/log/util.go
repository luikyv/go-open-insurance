package log

import (
	"context"
	"log/slog"
	"os"
)

const (
	CorrelationIDKey string = "correlation_id"
)

var baseLogger = slog.New(
	slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	),
)

func FromCtx(ctx context.Context) *slog.Logger {
	logger := baseLogger.With()

	correlationID := ctx.Value(CorrelationIDKey)
	if correlationID != nil {
		logger = logger.With(
			slog.String(CorrelationIDKey, correlationID.(string)),
		)
	}
	return logger
}
