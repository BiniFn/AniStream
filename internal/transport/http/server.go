package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/config"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Config *config.Env
	Router *chi.Mux
	Server *http.Server
}

func New(d *Dependencies) *App {
	r := chi.NewRouter()

	UseMiddlewares(d.Env, r)
	RegisterRoutes(r, d)

	srv := &http.Server{
		Addr:              ":" + d.Env.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &App{
		Config: d.Env,
		Router: r,
		Server: srv,
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
