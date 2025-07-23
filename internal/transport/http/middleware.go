package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/config"
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
		requestLogger(logger),
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

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return httplog.RequestLogger(logger, &httplog.Options{
		Level: slog.LevelInfo,
		LogExtraAttrs: func(req *http.Request, reqBody string, respStatus int) []slog.Attr {
			reqID := middleware.GetReqID(req.Context())
			return []slog.Attr{
				slog.String("request_id", reqID),
			}
		},
	})
}
