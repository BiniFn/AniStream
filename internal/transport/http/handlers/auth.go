package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/go-chi/chi/v5"
)

func MountAuthRoutes(r chi.Router, env *config.Env, userService *users.UserService) {
	r.Post("/login", login(env, userService))
	r.Get("/me", me(userService))
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

		http.SetCookie(w, &http.Cookie{
			Name:     "aniways_session",
			Value:    session.ID,
			Expires:  session.ExpiresAt.Time,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Domain:   env.CookieDomain,
		})

		jsonOK(w, user)
	}
}

func me(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("aniways_session")
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		user, err := userService.GetUserBySessionID(r.Context(), cookie.Value)
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		jsonOK(w, user)
	}
}
