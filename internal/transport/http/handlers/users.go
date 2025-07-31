package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/coeeter/aniways/internal/utils"
	"github.com/go-chi/chi/v5"
)

func MountUsersRoutes(r chi.Router, userService *users.UserService) {
	r.Post("/", createUser(userService))
	r.With(middleware.RequireUser).Group(func(r chi.Router) {
		r.Put("/", updateUser(userService))
		r.Delete("/", deleteUser(userService))
		r.Put("/password", updatePassword(userService))
		r.Put("/image", updateImage(userService))
	})
}

func createUser(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		user := &users.CreateUser{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			jsonError(w, http.StatusBadRequest, "Bad request")
			return
		}

		u, err := userService.CreateUser(r.Context(), user.Username, user.Email, user.Password)
		switch err {
		case nil:
			jsonOK(w, u)
		case users.ErrEmailTaken, users.ErrUsernameTaken:
			log.Warn("username or email already taken", "username", user.Username, "email", user.Email)
			jsonError(w, http.StatusConflict, err.Error())
		case users.ErrPasswordTooLong:
			log.Warn("password too long", "username", user.Username)
			jsonError(w, http.StatusBadRequest, err.Error())
		default:
			log.Error("user creation failed", "err", err)
			jsonError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func updatePassword(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user, _ := utils.CtxValue[users.User](r.Context())

		var updatePasswordBody struct {
			OldPassword string `json:"oldPassword"`
			NewPassword string `json:"newPassword"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updatePasswordBody); err != nil {
			jsonError(w, http.StatusBadRequest, "Bad request")
			return
		}

		err := userService.UpdatePassword(r.Context(), user.ID, updatePasswordBody.OldPassword, updatePasswordBody.NewPassword)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
		case users.ErrPasswordTooLong:
			log.Warn("password too long", "id", user.ID)
			jsonError(w, http.StatusBadRequest, err.Error())
		case users.ErrInvalidAuth:
			log.Warn("invalid password", "id", user.ID)
			jsonError(w, http.StatusUnauthorized, err.Error())
		default:
			log.Error("password update failed", "err", err)
			jsonError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func updateUser(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user, _ := utils.CtxValue[users.User](r.Context())

		var updateUserBody struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateUserBody); err != nil {
			jsonError(w, http.StatusBadRequest, "Bad request")
			return
		}

		u, err := userService.UpdateUser(r.Context(), user.ID, updateUserBody.Username, updateUserBody.Email)
		switch err {
		case nil:
			jsonOK(w, u)
		case users.ErrEmailTaken, users.ErrUsernameTaken:
			log.Warn("username or email already taken", "username", updateUserBody.Username, "email", updateUserBody.Email)
			jsonError(w, http.StatusConflict, err.Error())
		default:
			log.Error("user update failed", "err", err)
			jsonError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func deleteUser(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user, _ := utils.CtxValue[users.User](r.Context())

		var body struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, http.StatusBadRequest, "Bad request")
			return
		}

		err := userService.DeleteUser(r.Context(), user.ID, body.Password)
		switch err {
		case nil:
			http.SetCookie(w, &http.Cookie{
				Name:     "aniways_session",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  time.Now().Add(-time.Hour),
			})
			w.WriteHeader(http.StatusNoContent)
		case users.ErrInvalidAuth:
			log.Warn("invalid password", "id", user.ID)
			jsonError(w, http.StatusUnauthorized, err.Error())
		default:
			log.Error("user deletion failed", "err", err)
			jsonError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func updateImage(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user, _ := utils.CtxValue[users.User](r.Context())

		if err := r.ParseMultipartForm(10 << 20 /* 10 MiB */); err != nil {
			jsonError(w, http.StatusBadRequest, "bad request")
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			jsonError(w, http.StatusBadRequest, "image field is required")
			return
		}
		defer file.Close()

		err = userService.UpdateProfilePicture(r.Context(), user.ID, file)
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
		default:
			log.Error("user image update failed", "id", user.ID, "err", err)
			jsonError(w, http.StatusInternalServerError, "internal error")
		}
	}
}
