package handlers

import (
	"log"
	"net/http"

	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeEpisodesRoutes(r chi.Router, svc *animeSvc.AnimeService) {
	r.Route("/{id}/episodes", func(r chi.Router) {
		r.Get("/", getAnimeEpisodes(svc))
		r.Get("/{episodeID}/langs", getEpisodeLangs(svc))
		r.Get("/{episodeID}/stream/{type}", getEpisodeStream(svc))
		r.Get("/{episodeID}/stream/{type}/metadata", getEpisodeStreamMetadata(svc))
	})
}

func getAnimeEpisodes(svc *animeSvc.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
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

func getEpisodeLangs(svc *animeSvc.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
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

func getEpisodeStream(svc *animeSvc.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
			return
		}

		episodeID := chi.URLParam(r, "episodeID")
		if episodeID == "" {
			jsonError(w, http.StatusBadRequest, "episode ID is required")
			return
		}

		streamType := chi.URLParam(r, "type")
		if streamType == "" {
			jsonError(w, http.StatusBadRequest, "stream type is required")
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

func getEpisodeStreamMetadata(svc *animeSvc.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			jsonError(w, http.StatusBadRequest, "anime ID is required")
			return
		}

		episodeID := chi.URLParam(r, "episodeID")
		if episodeID == "" {
			jsonError(w, http.StatusBadRequest, "episode ID is required")
			return
		}

		streamType := chi.URLParam(r, "type")
		if streamType == "" {
			jsonError(w, http.StatusBadRequest, "stream type is required")
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
