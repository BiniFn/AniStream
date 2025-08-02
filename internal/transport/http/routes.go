package http

import (
	"net/http"

	"github.com/coeeter/aniways/internal/app"
	"github.com/coeeter/aniways/internal/transport/http/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, deps *app.Deps) {
	handler := handlers.New(deps, r)
	handler.RegisterRoutes()

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
