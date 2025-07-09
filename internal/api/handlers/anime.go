package handlers

import (
	"net/http"

	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeRoutes(r chi.Router, svc *animeSvc.Service) {
	r.Get("/{id}", getAnimeByID(svc))
	r.Get("/{id}/trailer", getAnimeTrailer(svc))
	r.Get("/{id}/episodes", getAnimeEpisodes(svc))
	r.Get("/genres", listGenres(svc))
	r.Get("/recently-updated", listRecentlyUpdated(svc))
}

func getAnimeByID(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
			return
		}

		resp, err := svc.GetAnimeByID(r.Context(), id)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
		jsonOK(w, resp)
	}
}

func listGenres(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.GetAnimeGenres(r.Context())
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime genres")
			return
		}
		jsonOK(w, resp)
	}
}

func listRecentlyUpdated(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetRecentlyUpdatedAnimes(r.Context(), page, size)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to fetch recently updated animes")
			return
		}
		jsonOK(w, resp)
	}
}

func getAnimeTrailer(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
			return
		}

		resp, err := svc.GetAnimeTrailer(r.Context(), id)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime trailer")
			return
		}
		jsonOK(w, resp)
	}
}

func getAnimeEpisodes(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
			return
		}

		resp, err := svc.GetAnimeEpisodes(r.Context(), id)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime episodes")
			return
		}
		jsonOK(w, resp)
	}
}
