package app

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger() *slog.Logger {
	var handler slog.Handler
	if os.Getenv("APP_ENV") == "development" {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	root := slog.New(handler)
	slog.SetDefault(root)
	return root
}
