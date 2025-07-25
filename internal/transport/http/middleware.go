package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/ctxutil"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

func UseMiddlewares(r *chi.Mux, logger *slog.Logger, d *Dependencies) {
	userService := users.NewUserService(d.Repo, d.Cld)
	r.Use(corsHandler(d.Env))

	r.Use(
		middleware.RealIP,
		middleware.RequestID,
		injectUser(userService),
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
			user, ok := ctxutil.Get[users.User](r.Context())
			if ok {
				logger = logger.With("user_id", user.ID)
			}
			ctx := ctxutil.Set(r.Context(), logger.With("request_id", reqID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func injectUser(userService *users.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("aniways_session")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := userService.GetUserBySessionID(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := ctxutil.Set(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, ok := ctxutil.Get[*slog.Logger](r.Context())
		if !ok {
			logger = slog.Default()
		}

		mw := httplog.RequestLogger(logger, &httplog.Options{
			Level: slog.LevelInfo,
		})

		mw(h).ServeHTTP(w, r)
	})
}
