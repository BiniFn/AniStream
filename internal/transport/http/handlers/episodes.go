package handlers

import (
	"log"
	"net/http"
	"strings"

	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeEpisodesRoutes(r chi.Router, svc *animeSvc.AnimeService) {
	r.Route("/{id}/episodes", func(r chi.Router) {
		r.Get("/", getAnimeEpisodes(svc))
		r.Get("/{episodeID}/servers", getServersForEpisode(svc))
		r.Get("/{episodeID}/langs", getEpisodeLangs(svc))
		r.Get("/{episodeID}/stream/{type}", getEpisodeStream(svc))
	})

	r.Get("/sources/{serverID}", getStreamingData(svc))
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

func getServersForEpisode(svc *animeSvc.AnimeService) http.HandlerFunc {
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

		resp, err := svc.GetServersForEpisode(r.Context(), id, episodeID)
		if err != nil {
			log.Printf("failed to fetch servers for episode %s of anime %s: %v", episodeID, id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch servers for episode")
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

func getStreamingData(svc *animeSvc.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverID := chi.URLParam(r, "serverID")
		if serverID == "" {
			jsonError(w, http.StatusBadRequest, "server ID is required")
			return
		}

		streamType := r.URL.Query().Get("type")
		if streamType == "" {
			jsonError(w, http.StatusBadRequest, "type is required")
			return
		}

		parts := strings.Split(streamType, "-")
		if len(parts) < 2 {
			jsonError(w, http.StatusBadRequest, "stream type must be in the format 'type-serverName'")
			return
		}

		streamType = parts[0]
		serverName := strings.Join(parts[1:], "-")

		resp, err := svc.GetStreamingData(r.Context(), serverID, streamType, serverName)
		if err != nil {
			log.Printf("failed to fetch streaming data for server %s: %v", serverID, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch streaming data")
			return
		}
		jsonOK(w, resp)
	}
}
