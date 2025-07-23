package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Config *config.Env
	Router *chi.Mux
	Server *http.Server
	Repo   *repository.Queries
	Cache  *cache.RedisClient
}

func New(config *config.Env, repo *repository.Queries, redis *cache.RedisClient) *App {
	r := chi.NewRouter()

	UseMiddlewares(config, r)
	RegisterRoutes(r, config, repo, redis)

	srv := &http.Server{
		Addr:              ":" + config.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &App{
		Router: r,
		Config: config,
		Server: srv,
		Repo:   repo,
		Cache:  redis,
	}
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		if err := a.Server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()
	log.Printf("🌐 AniWays API listening on http://localhost:%s", a.Config.AppPort)

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.Shutdown(shutCtx)
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	log.Println("🔻 server shut down gracefully")
	return nil
}
