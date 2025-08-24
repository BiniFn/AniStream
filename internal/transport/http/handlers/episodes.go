package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AnimeEpisodeRoutes() {
	h.r.Route("/anime/{id}/episodes", func(r chi.Router) {
		r.Get("/", h.getAnimeEpisodes)
		r.Get("/{episodeID}/servers", h.getEpisodeServers)
		r.Get("/servers/{serverID}", h.getEpisodeStreamData)
	})
}

func (h *Handler) getAnimeEpisodes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimeEpisodes(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

func (h *Handler) getEpisodeServers(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	episodeID, err := h.pathParam(r, "episodeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetEpisodeServers(r.Context(), id, episodeID)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

func (h *Handler) getEpisodeStreamData(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	serverID, err := h.pathParam(r, "serverID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	serverName := r.URL.Query().Get("server")
	if serverName == "" {
		h.jsonError(w, http.StatusBadRequest, "server name is required")
		return
	}
	streamType := r.URL.Query().Get("type")
	if streamType == "" {
		h.jsonError(w, http.StatusBadRequest, "stream type is required")
		return
	}

	resp, err := h.animeService.GetEpisodeStream(r.Context(), id, serverID, serverName, streamType)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch episode stream data", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch episode stream data")
		return
	}
}
