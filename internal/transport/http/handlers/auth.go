package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/service/auth"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AuthRoutes() {
	h.r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.login)
		r.Post("/forget-password", h.forgetPassword)
		r.Get("/u/{token}", h.getUser)
		r.Put("/reset-password/{token}", h.resetPassword)
		r.Post("/logout", h.logout)
		r.Get("/me", h.me)

		r.With(middleware.RequireUser).Group(func(r chi.Router) {
			r.Get("/providers", h.getProviders)
			r.Delete("/providers/{provider}", h.deleteProvider)
		})
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.userService.AuthenticateUser(r.Context(), body.Email, body.Password)
	if err != nil {
		log.Warn("Failed to authenticate user", "err", err)
		h.jsonError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	session, err := h.userService.CreateSession(r.Context(), user.ID)
	if err != nil {
		log.Error("Failed to create session", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	domain := "localhost"
	if h.deps.Env.CookieDomain != "" && h.deps.Env.CookieDomain != "localhost" {
		domain = fmt.Sprintf(".%s", h.deps.Env.CookieDomain) // enable subdomain cookies
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "aniways_session",
		Value:    session.ID,
		Expires:  session.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	})

	h.jsonOK(w, user)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.CtxValue[users.User](r.Context())
	if !ok {
		h.jsonError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	h.jsonOK(w, user)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	cookie, err := r.Cookie("aniways_session")
	if err != nil {
		h.jsonError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	err = h.userService.DeleteSession(r.Context(), cookie.Value)
	if err != nil {
		log.Error("Failed to delete session", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to delete session")
		return
	}

	domain := "localhost"
	if h.deps.Env.CookieDomain != "" && h.deps.Env.CookieDomain != "localhost" {
		domain = fmt.Sprintf(".%s", h.deps.Env.CookieDomain) // enable subdomain cookies
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "aniways_session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) forgetPassword(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	var input struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.authService.SendForgetPasswordEmail(r.Context(), input.Email); err != nil {
		log.Error("Failed to send reset password email", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to send reset password email")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	token, err := h.pathParam(r, "token")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	user, err := h.authService.GetUserByForgetPasswordToken(r.Context(), token)
	if err != nil {
		log.Error("Failed to get user by forget password token", "err", err)
		h.jsonError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	h.jsonOK(w, user)
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	token, err := h.pathParam(r, "token")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	var input struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	err = h.authService.ResetPassword(r.Context(), h.userService, token, input.Password)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case auth.ErrInvalidToken:
		h.jsonError(w, http.StatusUnauthorized, "Invalid token")
	case users.ErrPasswordTooLong:
		h.jsonError(w, http.StatusBadRequest, "Password is too long")
	default:
		log.Error("Failed to reset password", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to reset password")
	}
}

func (h *Handler) getProviders(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	user := middleware.GetUser(r)

	providers, err := h.authService.GetConnectedProviders(r.Context(), user.ID)
	if err != nil {
		log.Error("Failed to get providers", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to get providers")
		return
	}

	h.jsonOK(w, providers)
}

func (h *Handler) deleteProvider(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	provider, err := h.pathParam(r, "provider")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, "Invalid provider")
		return
	}

	err = h.authService.DisconnectProvider(r.Context(), user.ID, provider)
	if err != nil {
		log.Error("Failed to disconnect provider", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to disconnect provider")
		return
	}

	w.WriteHeader(http.StatusOK)
}
