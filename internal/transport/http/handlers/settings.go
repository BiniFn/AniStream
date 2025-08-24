package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/settings"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) SettingsRoutes() {
	h.r.With(middleware.RequireUser).Route("/settings", func(r chi.Router) {
		r.Get("/", h.getSettings)
		r.Post("/", h.saveSettings)
	})
}

func (h *Handler) getSettings(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	user := middleware.GetUser(r)

	settings, err := h.settingsService.GetSettings(r.Context(), user.ID)
	if err != nil {
		log.Error("failed to get settings", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get settings")
		return
	}

	h.jsonOK(w, settings)
}

func (h *Handler) saveSettings(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	user := middleware.GetUser(r)

	var req models.SettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	settings, err := h.settingsService.SaveSettings(r.Context(), settings.SaveSettingsParams{
		UserID:            user.ID,
		AutoNextEpisode:   req.AutoNextEpisode,
		AutoPlayEpisode:   req.AutoPlayEpisode,
		AutoResumeEpisode: req.AutoResumeEpisode,
		IncognitoMode:     req.IncognitoMode,
	})
	if err != nil {
		log.Error("failed to save settings", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to save settings")
		return
	}

	h.jsonOK(w, settings)
}
