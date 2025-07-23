package http

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/transport/http/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, deps *Dependencies) {
	refresher := anime.NewRefresher(deps.Repo, deps.MAL)
	animeService := anime.NewAnimeService(deps.Repo, refresher, deps.MAL, deps.Anilist, deps.Shiki, deps.Cache)
	userService := users.NewUserService(deps.Repo, deps.Cld)

	r.Route("/anime", func(r chi.Router) {
		handlers.MountAnimeRoutes(r, animeService)
		handlers.MountAnimeListingsRoutes(r, animeService)
		handlers.MountAnimeEpisodesRoutes(r, animeService)
	})

	r.Route("/users", func(r chi.Router) {
		handlers.MountUsersRoutes(r, userService)
	})

	r.Route("/auth", func(r chi.Router) {
		handlers.MountAuthRoutes(r, deps.Env, userService)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AniWays API"))
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
