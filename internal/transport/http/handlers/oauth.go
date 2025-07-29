package handlers

import (
	"fmt"
	"net/http"

	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func beginAuthHandler(provider oauth.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		state := uuid.NewString()

		url, err := provider.AuthURL(r.Context(), state)
		if err != nil {
			log.Error("unable to create oauth url", "provider", provider.Name(), "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to create auth url")
			return
		}

		log.Info("redirecting to oauth url", "provider", provider.Name(), "url", url)

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func callbackHandler(provider oauth.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		err := provider.ExchangeToken(r.Context(), user.ID, state, code)
		if err != nil {
			log.Error("unable to exchange token", "provider", provider.Name(), "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to exchange token")
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func mountOAuthRoutes(r chi.Router, provider oauth.Provider) {
	path := fmt.Sprintf("/oauth/%s", provider.Name())

	r.Route(path, func(r chi.Router) {
		r.Get("/", beginAuthHandler(provider))
		r.Get("/callback", callbackHandler(provider))
	})
}
