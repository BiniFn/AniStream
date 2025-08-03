package handlers

import (
	"encoding/json"
	"net/http"

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

	var body struct {
		AutoNextEpisode   bool `json:"autoNextEpisode"`
		AutoPlayEpisode   bool `json:"autoPlayEpisode"`
		AutoResumeEpisode bool `json:"autoResumeEpisode"`
		IncognitoMode     bool `json:"incognitoMode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	settings, err := h.settingsService.SaveSettings(r.Context(), settings.SaveSettingsParams{
		UserID:            user.ID,
		AutoNextEpisode:   body.AutoNextEpisode,
		AutoPlayEpisode:   body.AutoPlayEpisode,
		AutoResumeEpisode: body.AutoResumeEpisode,
		IncognitoMode:     body.IncognitoMode,
	})
	if err != nil {
		log.Error("failed to save settings", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to save settings")
		return
	}

	h.jsonOK(w, settings)
}
