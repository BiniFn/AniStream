package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/ctxutil"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/go-chi/chi/v5"
)

func MountAuthRoutes(r chi.Router, env *config.Env, userService *users.UserService) {
	r.Post("/login", login(env, userService))
	r.Post("/logout", logout(env, userService))
	r.Get("/me", me)
}

func login(env *config.Env, userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		user, err := userService.AuthenticateUser(r.Context(), body.Email, body.Password)
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		session, err := userService.CreateSession(r.Context(), user.ID)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to create token")
			return
		}

		domain := "localhost"
		if env.CookieDomain != "" && env.CookieDomain != "localhost" {
			domain = fmt.Sprintf(".%s", env.CookieDomain) // enable subdomain cookies
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

		jsonOK(w, user)
	}
}

func me(w http.ResponseWriter, r *http.Request) {
	user, ok := ctxutil.Get[users.User](r.Context())
	if !ok {
		jsonError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	jsonOK(w, user)
}

func logout(env *config.Env, userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("aniways_session")
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		err = userService.DeleteSession(r.Context(), cookie.Value)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "Failed to delete session")
			return
		}

		domain := "localhost"
		if env.CookieDomain != "" && env.CookieDomain != "localhost" {
			domain = fmt.Sprintf(".%s", env.CookieDomain) // enable subdomain cookies
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
}
