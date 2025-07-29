package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/config"
	"github.com/coeeter/aniways/internal/ctxutil"
	"github.com/coeeter/aniways/internal/service/auth"
	"github.com/coeeter/aniways/internal/service/auth/oauth"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAuthRoutes(
	r chi.Router,
	env *config.Env,
	userService *users.UserService,
	authService *auth.AuthService,
	redis *cache.RedisClient,
) {
	r.Post("/login", login(env, userService))
	r.Post("/forget-password", forgetPassword(authService))
	r.Get("/u/{token}", getUser(authService))
	r.Put("/reset-password/{token}", resetPassword(authService, userService))
	r.Post("/logout", logout(env, userService))
	r.Get("/me", me)

	r.With(middleware.RequireUser).Group(func(r chi.Router) {
		mal := oauth.NewMALProvider(
			env.MyAnimeListClientID,
			env.MyAnimeListClientSecret,
			fmt.Sprintf("%s/auth/oauth/myanimelist/callback", env.ApiURL),
			redis,
		)

		MountOAuthRoutes(r, mal)
	})
}

func login(env *config.Env, userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

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
			log.Warn("Failed to authenticate user", "err", err)
			jsonError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		session, err := userService.CreateSession(r.Context(), user.ID)
		if err != nil {
			log.Error("Failed to create session", "err", err)
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
		log := logger(r)

		cookie, err := r.Cookie("aniways_session")
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		err = userService.DeleteSession(r.Context(), cookie.Value)
		if err != nil {
			log.Error("Failed to delete session", "err", err)
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

func forgetPassword(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		var input struct {
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request")
			return
		}

		if err := authService.SendForgetPasswordEmail(r.Context(), input.Email); err != nil {
			log.Error("Failed to send reset password email", "err", err)
			jsonError(w, http.StatusInternalServerError, "Failed to send reset password email")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getUser(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		token, err := pathParam(r, "token")
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid token")
			return
		}

		user, err := authService.GetUserByForgetPasswordToken(r.Context(), token)
		if err != nil {
			log.Error("Failed to get user by forget password token", "err", err)
			jsonError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		jsonOK(w, user)
	}
}

func resetPassword(authService *auth.AuthService, userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		token, err := pathParam(r, "token")
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid token")
			return
		}

		var input struct {
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid request")
			return
		}

		err = authService.ResetPassword(r.Context(), userService, token, input.Password)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
		case auth.ErrInvalidToken:
			jsonError(w, http.StatusUnauthorized, "Invalid token")
		case users.ErrPasswordTooLong:
			jsonError(w, http.StatusBadRequest, "Password is too long")
		default:
			log.Error("Failed to reset password", "err", err)
			jsonError(w, http.StatusInternalServerError, "Failed to reset password")
		}
	}
}
