package handlers

import (
	"log"
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
		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimeEpisodes(r.Context(), id)
		if err != nil {
			log.Printf("failed to fetch anime episodes for ID %s: %v", id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime episodes")
			return
		}
		jsonOK(w, resp)
	}
}

func getEpisodeLangs(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			log.Printf("failed to fetch languages for episode %s of anime %s: %v", episodeID, id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch languages for episode")
			return
		}
		jsonOK(w, resp)
	}
}

func getEpisodeStream(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		streamType, err := pathParam(r, "streamType")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetEpisodeStream(r.Context(), id, episodeID, streamType)
		if err != nil {
			log.Printf("failed to fetch stream for episode %s of anime %s: %v", episodeID, id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch stream for episode")
			return
		}
		jsonOK(w, resp)
	}
}

func getEpisodeStreamMetadata(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		streamType, err := pathParam(r, "streamType")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetStreamMetadata(r.Context(), id, episodeID, streamType)
		if err != nil {
			log.Printf("failed to fetch stream metadata for episode %s of anime %s: %v", episodeID, id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch stream metadata for episode")
			return
		}
		jsonOK(w, resp)
	}
}
