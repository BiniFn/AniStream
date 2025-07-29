package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/database"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/transport/http"
	"github.com/coeeter/aniways/internal/worker"
	"github.com/lmittmann/tint"
)

func newRootLogger() *slog.Logger {
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

func main() {
	env, err := config.LoadEnv()
	rootLog := newRootLogger()
	if err != nil {
		rootLog.Error("Error loading environment variables:", "err", err)
		os.Exit(1)
	}

	dbLog := rootLog.With("component", "database")
	db, err := database.New(env, dbLog)
	if err != nil {
		rootLog.Error("Error connecting to the database:", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	redisLog := rootLog.With("component", "redis")
	redis, err := cache.NewRedisClient(context.Background(), env.AppEnv, env.RedisAddr, env.RedisPassword, redisLog)
	if err != nil {
		rootLog.Error("Error connecting to Redis:", "err", err)
		os.Exit(1)
	}
	defer redis.Close()

	repo := repository.New(db)
	scraper := hianime.NewHianimeScraper()

	deps, err := http.BuildDeps(env, repo, redis)
	if err != nil {
		rootLog.Error("failed to build dependencies:", "err", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	workerLog := rootLog.With("component", "worker")
	mgr := worker.NewManager(repo, scraper, redis, workerLog)

	if err := mgr.Bootstrap(ctx); err != nil {
		rootLog.Error("Error in bootstrapping:", "err", err)
		os.Exit(1)
	}

	mgr.StartBackground(ctx, deps.Providers)

	httpLog := rootLog.With("component", "http")
	app := http.New(deps, httpLog)
	if err := app.Run(ctx); err != nil {
		rootLog.Error("failed to run application", "err", err)
		os.Exit(1)
	}
}
