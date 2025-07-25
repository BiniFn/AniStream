package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/config"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Config *config.Env
	Router *chi.Mux
	Server *http.Server
	Log    *slog.Logger
}

func New(d *Dependencies, log *slog.Logger) *App {
	r := chi.NewRouter()

	UseMiddlewares(r, log, d)
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
		Log:    log,
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
	a.Log.Info("AniWays API listening", "on", a.Config.AppPort)

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
	a.Log.Info("server shut down gracefully")
	return nil
}
