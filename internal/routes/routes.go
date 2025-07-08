package routes

import (
	"net/http"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/go-chi/chi/v5"
)

func MountGlobal(r *chi.Mux, repo *repository.Queries) {
	mountAnimeRoutes(r, repo)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("AniWays API"))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
