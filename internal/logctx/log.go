package logctx

import (
	"context"
	"log/slog"
)

type logkey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, logkey{}, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(logkey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
