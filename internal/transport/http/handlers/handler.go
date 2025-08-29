package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/coeeter/aniways/docs"
	"github.com/coeeter/aniways/internal/app"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/coeeter/aniways/internal/service/auth"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/coeeter/aniways/internal/service/library"
	"github.com/coeeter/aniways/internal/service/settings"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	r               *chi.Mux
	deps            *app.Deps
	validator       *validator.Validate
	refresher       *anime.MetadataRefresher
	animeService    *anime.AnimeService
	userService     *users.UserService
	authService     *auth.AuthService
	settingsService *settings.SettingsService
	libraryService  *library.LibraryService
	providerMap     map[string]oauth.Provider
}

func New(deps *app.Deps, r *chi.Mux) *Handler {
	refresher := anime.NewRefresher(deps.Repo, deps.MAL)
	animeService := anime.NewAnimeService(deps.Repo, refresher, deps.MAL, deps.Anilist, deps.Shiki, deps.Cache)
	userService := users.NewUserService(deps.Repo, deps.Cld)
	authService := auth.NewAuthService(deps.Repo, deps.EmailClient, deps.Env.FrontendURL)
	settingsService := settings.NewSettingsService(deps.Repo)
	libraryService := library.NewLibraryService(deps.Repo, refresher)

	providerMap := make(map[string]oauth.Provider)
	for _, provider := range deps.Providers {
		providerMap[provider.Name()] = provider
	}

	return &Handler{
		r:               r,
		deps:            deps,
		validator:       validator.New(),
		refresher:       refresher,
		animeService:    animeService,
		userService:     userService,
		authService:     authService,
		settingsService: settingsService,
		libraryService:  libraryService,
		providerMap:     providerMap,
	}
}

func (h *Handler) RegisterRoutes() {
	h.r.Get("/", h.home)
	h.r.Get("/healthz", h.healthz)

	h.AnimeDetailsRoutes()
	h.AnimeListingRoutes()
	h.AnimeEpisodeRoutes()
	h.AuthRoutes()
	h.OauthRoutes()
	h.UserRoutes()
	h.LibraryRoutes()
	h.SettingsRoutes()

	h.RegisterOpenAPIRoutes()

	h.r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}

func (h *Handler) RegisterOpenAPIRoutes() {
	if h.deps.Env.AppEnv == "development" {
		h.r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.Write(docs.OpenAPISpec)
		})

		h.r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(docs.OpenAPISpec)))
	}
}

func (h *Handler) home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("AniWays API\n"))
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Router() *chi.Mux {
	return h.r
}

func (h *Handler) parsePagination(r *http.Request, defaultPage, defaultSize int) (page, size int, err error) {
	q := r.URL.Query()

	page = defaultPage
	if v := q.Get("page"); v != "" {
		pageVal, err2 := strconv.Atoi(v)
		if err2 != nil || pageVal < 1 {
			return 0, 0, fmt.Errorf("invalid page")
		}
		page = pageVal
	}

	size = defaultSize
	if v := q.Get("itemsPerPage"); v != "" {
		sizeVal, err2 := strconv.Atoi(v)
		if err2 != nil || sizeVal < 1 {
			return 0, 0, fmt.Errorf("invalid itemsPerPage")
		}
		size = sizeVal
	}

	return page, size, nil
}

func (h *Handler) pathParam(r *http.Request, key string) (string, error) {
	v := chi.URLParam(r, key)
	if v == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}

func (h *Handler) jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}

func (h *Handler) jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) logger(r *http.Request) *slog.Logger {
	logger, ok := utils.CtxValue[*slog.Logger](r.Context())
	if !ok {
		return slog.Default()
	}
	return logger.With("layer", "controller")
}

func (h *Handler) parseAndValidate(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return false
	}

	if err := h.validator.Struct(req); err != nil {
		h.jsonValidationError(w, err)
		return false
	}

	return true
}

func (h *Handler) jsonValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Validation failed"})
		return
	}

	details := make(map[string]string)
	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()
		switch fieldErr.Tag() {
		case "required":
			details[field] = field + " is required"
		case "email":
			details[field] = "Invalid email format"
		case "min":
			details[field] = field + " must be at least " + fieldErr.Param() + " characters"
		case "max":
			details[field] = field + " must be at most " + fieldErr.Param() + " characters"
		case "oneof":
			details[field] = field + " must be one of: " + fieldErr.Param()
		default:
			details[field] = field + " is invalid"
		}
	}

	_ = json.NewEncoder(w).Encode(models.ValidationErrorResponse{
		Error:   "Validation failed",
		Details: details,
	})
}

