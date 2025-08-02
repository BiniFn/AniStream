package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AnimeEpisodeRoutes() {
	h.r.Route("/anime/{id}/episodes", func(r chi.Router) {
		r.Get("/", h.getAnimeEpisodes)
		r.Get("/{episodeID}/langs", h.getEpisodeLangs)
		r.Get("/{episodeID}/stream/{type}", h.getEpisodeStream)
		r.Get("/{episodeID}/stream/{type}/metadata", h.getEpisodeStreamMetadata)
	})
}

// getAnimeEpisodes returns the episodes of an anime
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

// getEpisodeLangs returns the languages of an episode
func (h *Handler) getEpisodeLangs(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	episodeID := chi.URLParam(r, "episodeID")
	if episodeID == "" {
		h.jsonError(w, http.StatusBadRequest, "episode ID is required")
		return
	}

	resp, err := h.animeService.GetEpisodeLangs(r.Context(), id, episodeID)
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

// getEpisodeStream returns the stream of an episode
func (h *Handler) getEpisodeStream(w http.ResponseWriter, r *http.Request) {
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

	streamType, err := h.pathParam(r, "type")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetEpisodeStream(r.Context(), id, episodeID, streamType)
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

// getEpisodeStreamMetadata returns the metadata of an episode stream
func (h *Handler) getEpisodeStreamMetadata(w http.ResponseWriter, r *http.Request) {
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

	streamType, err := h.pathParam(r, "type")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetStreamMetadata(r.Context(), id, episodeID, streamType)
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
