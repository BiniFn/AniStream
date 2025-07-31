package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/coeeter/aniways/internal/app"
	"github.com/coeeter/aniways/internal/worker"
)

func main() {
	deps, err := app.InitDeps(context.Background(), "WORKER")
	if err != nil {
		deps.Log.Error("Error initializing dependencies:", "err", err)
		os.Exit(1)
	}
	defer deps.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mgr := worker.NewManager(deps.Repo, deps.Scraper, deps.Cache, deps.Log.With("component", "worker"))

	if err := mgr.Bootstrap(ctx); err != nil {
		deps.Log.Error("Error in bootstrapping:", "err", err)
		os.Exit(1)
	}

	mgr.StartBackground(ctx, deps.Providers)

	<-ctx.Done()
}
