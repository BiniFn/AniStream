package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type GetByIDInput struct {
	ID int32 `in:"path=id"`
}

func (h *Handler) CharacterRoutes() {
	h.r.With(httpin.NewInput(GetByIDInput{})).Route("/characters/{id}", func(r chi.Router) {
		r.Get("/", h.getCharacterByID)
	})

	h.r.With(httpin.NewInput(GetByIDInput{})).Route("/characters/va/{id}", func(r chi.Router) {
		r.Get("/", h.getPersonByID)
	})
}

// @Summary Get character by ID
// @Description Get detailed character information by MAL ID
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path int32 true "Character MAL ID"
// @Success 200 {object} models.CharacterFullResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id} [get]
func (h *Handler) getCharacterByID(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetByIDInput)
	id := input.ID

	resp, err := h.animeService.GetCharacterFull(r.Context(), id)
	switch err {
	case anime.ErrCharacterNotFound:
		log.Warn("character not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "character not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch character details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch character details")
		return
	}
}

// @Summary Get person by ID
// @Description Get detailed person (voice actor) information by MAL ID
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path int32 true "Person MAL ID"
// @Success 200 {object} models.PersonFullResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/va/{id} [get]
func (h *Handler) getPersonByID(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetByIDInput)
	id := input.ID

	resp, err := h.animeService.GetPersonFull(r.Context(), id)
	switch err {
	case anime.ErrPersonNotFound:
		log.Warn("person not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "person not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch person details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch person details")
		return
	}
}

