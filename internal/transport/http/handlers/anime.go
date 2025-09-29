package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type GetAnimeByIDInput struct {
	ID string `in:"path=id"`
}

func (h *Handler) AnimeDetailsRoutes() {
	h.r.With(httpin.NewInput(GetAnimeByIDInput{})).Route("/anime/{id}", func(r chi.Router) {
		r.Get("/", h.getAnimeByID)
		r.Get("/trailer", h.getAnimeTrailer)
		r.Get("/banner", h.getAnimeBanner)
		r.Get("/franchise", h.getAnimeFranchise)
	})
}

// @Summary Get anime by ID
// @Description Get anime by ID
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.AnimeWithMetadataResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id} [get]
func (h *Handler) getAnimeByID(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

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

// @Summary Get anime trailer
// @Description Get anime trailer
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.TrailerResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/trailer [get]
func (h *Handler) getAnimeTrailer(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.animeService.GetAnimeTrailer(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case anime.TrailerNotFound:
		log.Warn("trailer not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "trailer not found")
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

// @Summary Get anime banner
// @Description Get anime banner
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.BannerResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/banner [get]
func (h *Handler) getAnimeBanner(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.animeService.GetAnimeBanner(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case anime.BannerNotFound:
		log.Warn("banner not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "banner not found")
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

// @Summary Get anime franchise relations
// @Description Get anime franchise relations
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.RelationsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/franchise [get]
func (h *Handler) getAnimeFranchise(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

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
