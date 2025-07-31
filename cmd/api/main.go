package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/coeeter/aniways/internal/app"
	"github.com/coeeter/aniways/internal/transport/http"
)

func main() {
	deps, err := app.InitDeps(context.Background(), "API")
	if err != nil {
		deps.Log.Error("Error initializing dependencies:", "err", err)
		os.Exit(1)
	}
	defer deps.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	httpLog := deps.Log.With("component", "http")
	app := http.New(deps, httpLog)

	if err := app.Run(ctx); err != nil {
		deps.Log.Error("failed to run application", "err", err)
		os.Exit(1)
	}
}
