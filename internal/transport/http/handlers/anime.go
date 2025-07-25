package handlers

import (
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
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeByID(r.Context(), id)
		switch err {
		case anime.ErrAnimeNotFound:
			log.Warn("anime not found", "id", id, "err", err)
			jsonError(w, http.StatusNotFound, "anime not found")
			return
		case nil:
			jsonOK(w, resp)
			return
		default:
			log.Error("failed to fetch anime details", "id", id, "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
	}
}

func getAnimeTrailer(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeTrailer(r.Context(), id)
		switch err {
		case anime.ErrAnimeNotFound:
			log.Warn("anime not found", "id", id, "err", err)
			jsonError(w, http.StatusNotFound, "anime not found")
			return
		case nil:
			jsonOK(w, resp)
			return
		default:
			log.Error("failed to fetch anime details", "id", id, "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
	}
}

func getAnimeBanner(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeBanner(r.Context(), id)
		switch err {
		case anime.ErrAnimeNotFound:
			log.Warn("anime not found", "id", id, "err", err)
			jsonError(w, http.StatusNotFound, "anime not found")
			return
		case nil:
			jsonOK(w, resp)
			return
		default:
			log.Error("failed to fetch anime details", "id", id, "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
	}
}

func franchise(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeRelations(r.Context(), id)
		switch err {
		case anime.ErrAnimeNotFound:
			log.Warn("anime not found", "id", id, "err", err)
			jsonError(w, http.StatusNotFound, "anime not found")
			return
		case nil:
			jsonOK(w, resp)
			return
		default:
			log.Error("failed to fetch anime details", "id", id, "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
	}
}
