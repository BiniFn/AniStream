package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeEpisodesRoutes(r chi.Router, svc *anime.AnimeService) {
	r.Route("/{id}/episodes", func(r chi.Router) {
		r.Get("/", getAnimeEpisodes(svc))
		r.Get("/{episodeID}/langs", getEpisodeLangs(svc))
		r.Get("/{episodeID}/stream/{type}", getEpisodeStream(svc))
		r.Get("/{episodeID}/stream/{type}/metadata", getEpisodeStreamMetadata(svc))
	})
}

func getAnimeEpisodes(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeEpisodes(r.Context(), id)
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

func getEpisodeLangs(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		episodeID := chi.URLParam(r, "episodeID")
		if episodeID == "" {
			jsonError(w, http.StatusBadRequest, "episode ID is required")
			return
		}

		resp, err := svc.GetEpisodeLangs(r.Context(), id, episodeID)
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

func getEpisodeStream(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		episodeID, err := pathParam(r, "episodeID")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		streamType, err := pathParam(r, "type")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetEpisodeStream(r.Context(), id, episodeID, streamType)
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

func getEpisodeStreamMetadata(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		episodeID, err := pathParam(r, "episodeID")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		streamType, err := pathParam(r, "type")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetStreamMetadata(r.Context(), id, episodeID, streamType)
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
