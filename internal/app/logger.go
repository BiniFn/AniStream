package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/lmittmann/tint"
)

func NewLogger(svcName string) *slog.Logger {
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

	color := "blue"
	switch svcName {
	case "WORKER":
		color = "green"
	case "API":
		color = "red"
	}
	txt := fmt.Sprintf("AniStream %s", svcName)
	fig := figure.NewColorFigure(txt, "", color, true)
	fig.Print()

	return root
}
