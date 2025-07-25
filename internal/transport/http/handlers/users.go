package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coeeter/aniways/internal/service/users"
	"github.com/go-chi/chi/v5"
)

func MountUsersRoutes(r chi.Router, userService *users.UserService) {
	r.Post("/", createUser(userService))
}

func createUser(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		user := &users.CreateUser{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
