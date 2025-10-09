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

// @Summary Get anime episodes
// @Description Get anime episodes
// @Tags Episodes
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.EpisodeListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/episodes [get]
func (h *Handler) getAnimeEpisodes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.Anime.GetAnimeEpisodes(r.Context(), id)
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

// @Summary Get episode servers
// @Description Get episode servers
// @Tags Episodes
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Param episodeID path string true "Episode ID"
// @Success 200 {object} models.EpisodeServerListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/episodes/{episodeID}/servers [get]
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

	resp, err := h.services.Anime.GetEpisodeServers(r.Context(), id, episodeID)
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

// @Summary Get episode stream data
// @Description Get episode stream data
// @Tags Episodes
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Param serverID path string true "Server ID"
// @Param server query string true "Server name"
// @Param type query string true "Stream type"
// @Success 200 {object} models.StreamingDataResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/episodes/servers/{serverID} [get]
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

	resp, err := h.services.Anime.GetEpisodeStream(r.Context(), id, serverID, serverName, streamType)
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
