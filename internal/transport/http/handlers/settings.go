package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coeeter/aniways/internal/service/settings"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func MountSettingsRoute(r chi.Router, svc *settings.SettingsService) {
	r.With(middleware.RequireUser).Group(func(r chi.Router) {
		r.Get("/", getSettings(svc))
		r.Post("/", saveSettings(svc))
	})
}

func getSettings(svc *settings.SettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		user := middleware.GetUser(r)

		settings, err := svc.GetSettings(r.Context(), user.ID)
		if err != nil {
			log.Error("failed to get settings", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get settings")
			return
		}

		jsonOK(w, settings)
	}
}

func saveSettings(svc *settings.SettingsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		user := middleware.GetUser(r)

		var body struct {
			AutoNextEpisode   bool `json:"autoNextEpisode"`
			AutoPlayEpisode   bool `json:"autoPlayEpisode"`
			AutoResumeEpisode bool `json:"autoResumeEpisode"`
			IncognitoMode     bool `json:"incognitoMode"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		settings, err := svc.SaveSettings(r.Context(), settings.SaveSettingsParams{
			UserID:            user.ID,
			AutoNextEpisode:   body.AutoNextEpisode,
			AutoPlayEpisode:   body.AutoPlayEpisode,
			AutoResumeEpisode: body.AutoResumeEpisode,
			IncognitoMode:     body.IncognitoMode,
		})
		if err != nil {
			log.Error("failed to save settings", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to save settings")
			return
		}

		jsonOK(w, settings)
	}
}
