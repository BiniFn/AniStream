package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	r.Use(corsHandler(config))

	srv := &http.Server{
		Addr:              ":" + config.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
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

func (a *App) Run() error {
	// wire routes…
	MountGlobalRoutes(a.Router, a.Config, a.Repo, a.Cache)

	// start listening
	errChan := make(chan error, 1)
	go func() {
		errChan <- a.Server.ListenAndServe()
	}()
	log.Printf("🌐 AniWays API listening on http://localhost:%s", a.Config.AppPort)

	// wait for either a server error or an OS signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-stop:
		log.Printf("🛑 received %v, shutting down…", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.Shutdown(ctx)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	log.Println("🔻 server shut down gracefully")
	return nil
}

func corsHandler(env *config.Env) func(http.Handler) http.Handler {
	if env.AppEnv == "development" {
		return cors.AllowAll().Handler
	}

	// In production, use specific allowed origins
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.AllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           300, // 5 minutes
		AllowCredentials: true,
	})
}
