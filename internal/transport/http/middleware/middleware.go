package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

type MiddlewareConfig struct {
	Router *chi.Mux
	Logger *slog.Logger
	Env    *config.Env
	Repo   *repository.Queries
	Cld    *cloudinary.Cloudinary
}

func UseMiddlewares(c MiddlewareConfig) {
	userService := users.NewUserService(c.Repo, c.Cld)
	c.Router.Use(corsHandler(c.Env))

	c.Router.Use(
		middleware.RealIP,
		middleware.RequestID,
		injectUser(userService),
		injectLogger(c.Logger),
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
			reqLogger := logger.With("request_id", middleware.GetReqID(r.Context()))

			if user, ok := utils.CtxValue[models.UserResponse](r.Context()); ok {
				reqLogger = reqLogger.With("user_id", user.ID)
			}

			ctx := utils.CtxWithValue(r.Context(), reqLogger)
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

			ctx := utils.CtxWithValue(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func requestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, ok := utils.CtxValue[*slog.Logger](r.Context())
		if !ok {
			logger = slog.Default()
		}

		mw := httplog.RequestLogger(logger, &httplog.Options{
			Level: slog.LevelInfo,
		})

		mw(h).ServeHTTP(w, r)
	})
}
