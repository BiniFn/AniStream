package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) OauthRoutes() {
	h.r.With(middleware.RequireUser).Route("/auth/oauth", func(r chi.Router) {
		r.Get("/{provider}", h.beginAuthHandler)
		r.Get("/{provider}/callback", h.callbackHandler)
	})
}

func (h *Handler) beginAuthHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	state := uuid.NewString()

	provider, ok := h.providerMap[chi.URLParam(r, "provider")]
	if !ok {
		h.jsonError(w, http.StatusNotFound, "provider not found")
		return
	}

	url, err := provider.AuthURL(r.Context(), state)
	if err != nil {
		log.Error("unable to create oauth url", "provider", provider.Name(), "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to create auth url")
		return
	}

	log.Info("redirecting to oauth url", "provider", provider.Name(), "url", url)

	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) callbackHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	provider, ok := h.providerMap[chi.URLParam(r, "provider")]
	if !ok {
		h.jsonError(w, http.StatusNotFound, "provider not found")
		return
	}

	err := provider.ExchangeToken(r.Context(), user.ID, state, code)
	if err != nil {
		log.Error("unable to exchange token", "provider", provider.Name(), "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to exchange token")
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
