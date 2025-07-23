package handlers

import (
	"log"
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeRoutes(r chi.Router, svc *anime.AnimeService) {
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getAnimeByID(svc))
		r.Get("/trailer", getAnimeTrailer(svc))
		r.Get("/banner", getAnimeBanner(svc))
		r.Get("/franchise", franchise(svc))
	})
}

func getAnimeByID(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeByID(r.Context(), id)
		if err != nil {
			log.Printf("failed to fetch anime details for ID %s: %v", id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
		jsonOK(w, resp)
	}
}

func getAnimeTrailer(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeTrailer(r.Context(), id)
		if err != nil {
			log.Printf("failed to fetch trailer for anime %s: %v", id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime trailer")
			return
		}
		jsonOK(w, resp)
	}
}

func getAnimeBanner(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeBanner(r.Context(), id)
		if err != nil {
			log.Printf("failed to fetch banner for anime %s: %v", id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime banner")
			return
		}
		jsonOK(w, resp)
	}
}

func franchise(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeRelations(r.Context(), id)
		if err != nil {
			log.Printf("failed to fetch franchise for anime ID %s: %v", id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch franchise")
			return
		}
		jsonOK(w, resp)
	}
}
