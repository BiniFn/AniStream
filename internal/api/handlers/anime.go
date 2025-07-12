package handlers

import (
	"log"
	"net/http"
	"strings"

	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeRoutes(r chi.Router, svc *animeSvc.Service) {
	r.Get("/{id}", getAnimeByID(svc))
	r.Get("/{id}/trailer", getAnimeTrailer(svc))
	r.Get("/{id}/episodes", getAnimeEpisodes(svc))
	r.Get("/{id}/episodes/{episodeID}", getServersForEpisode(svc))
	r.Get("/sources/{serverID}", getStreamingData(svc))
	r.Get("/genres", listGenres(svc))
	r.Get("/recently-updated", listRecentlyUpdated(svc))
	r.Get("/search", searchAnimes(svc))
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
			log.Printf("failed to fetch anime details for ID %s: %v", id, err)
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
			log.Printf("failed to fetch anime genres: %v", err)
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
			log.Printf("failed to fetch recently updated animes: %v", err)
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
			log.Printf("failed to fetch trailer for anime %s: %v", id, err)
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
			log.Printf("failed to fetch anime episodes for ID %s: %v", id, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime episodes")
			return
		}
		jsonOK(w, resp)
	}
}

func searchAnimes(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			jsonError(w, http.StatusBadRequest, "search query is required")
			return
		}
		if len(query) < 3 {
			jsonError(w, http.StatusBadRequest, "search query must be at least 3 characters long")
			return
		}

		genre := r.URL.Query().Get("genre")
		if genre != "" && len(genre) < 3 {
			jsonError(w, http.StatusBadRequest, "genre query must be at least 3 characters long")
			return
		}

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.SearchAnimes(r.Context(), query, genre, page, size)
		if err != nil {
			log.Printf("failed to search animes: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to search animes")
			return
		}
		jsonOK(w, resp)
	}
}

func getServersForEpisode(svc *animeSvc.Service) http.HandlerFunc {
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

func getStreamingData(svc *animeSvc.Service) http.HandlerFunc {
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
