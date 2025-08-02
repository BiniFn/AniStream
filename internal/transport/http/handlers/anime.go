package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AnimeDetailsRoutes() {
	h.r.Route("/anime/{id}", func(r chi.Router) {
		r.Get("/", h.getAnimeByID)
		r.Get("/trailer", h.getAnimeTrailer)
		r.Get("/banner", h.getAnimeBanner)
		r.Get("/franchise", h.getAnimeFranchise)
	})
}

// getAnimeByID retrieves an anime by its ID.
func (h *Handler) getAnimeByID(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimeByID(r.Context(), id)
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

// getAnimeTrailer retrieves the trailer of an anime by its ID.
func (h *Handler) getAnimeTrailer(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimeTrailer(r.Context(), id)
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

// getAnimeBanner retrieves the banner of an anime by its ID.
func (h *Handler) getAnimeBanner(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimeBanner(r.Context(), id)
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

// getAnimeFranchise retrieves the franchise of an anime by its ID.
func (h *Handler) getAnimeFranchise(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimeRelations(r.Context(), id)
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
