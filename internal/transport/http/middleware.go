package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/logctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

func UseMiddlewares(config *config.Env, r *chi.Mux, logger *slog.Logger) {
	r.Use(corsHandler(config))

	r.Use(
		middleware.RealIP,
		middleware.RequestID,
		injectLogger(logger),
		requestLogger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)
}

func corsHandler(env *config.Env) func(http.Handler) http.Handler {
	if env.AppEnv == "development" {
		return cors.AllowAll().Handler
	}

	// In production, use specific allowed origins
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.AllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		MaxAge:           300, // 5 minutes
		AllowCredentials: true,
	})
}

func injectLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ctx := logctx.WithLogger(r.Context(), logger.With("request_id", reqID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logctx.Logger(r.Context())

		mw := httplog.RequestLogger(logger, &httplog.Options{
			Level: slog.LevelInfo,
		})

		mw(h).ServeHTTP(w, r)
	})
}
