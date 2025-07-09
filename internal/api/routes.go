package api

import (
	"net/http"

	"github.com/coeeter/aniways/internal/api/handlers"
	"github.com/coeeter/aniways/internal/repository"
	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountGlobalRoutes(r *chi.Mux, repo *repository.Queries) {
	r.Route("/anime", func(r chi.Router) {
		svc := animeSvc.New(repo)
		handlers.MountAnimeRoutes(r, svc)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AniWays API"))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
