package main

import (
	"context"
	"log"
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
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	db, err := database.New(env)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	redis, err := cache.NewRedisClient(context.Background(), env.AppEnv, env.RedisAddr, env.RedisPassword)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	defer redis.Close()

	repo := repository.New(db)
	scraper := hianime.NewHianimeScraper()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	m := worker.NewManager(repo, scraper, redis)

	if err := m.Bootstrap(ctx); err != nil {
		log.Fatalf("Error in bootstrapping: %v", err)
	}

	m.StartBackground(ctx)

	deps, err := http.BuildDeps(env, repo, redis)
	if err != nil {
		log.Fatalf("failed to build dependencies: %v", err)
	}

	app := http.New(deps)
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
