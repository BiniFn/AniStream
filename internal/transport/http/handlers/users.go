package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coeeter/aniways/internal/service/users"
	"github.com/go-chi/chi/v5"
)

func MountUsersRoutes(r chi.Router, userService *users.UserService) {
	r.Post("/", createUser(userService))
}

func createUser(userService *users.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &users.CreateUser{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := userService.CreateUser(r.Context(), user.Username, user.Email, user.Password)

		if err != nil {
			log.Printf("Error creating user: %v", err)
			switch err {
			case users.ErrEmailTaken, users.ErrUsernameTaken:
				jsonError(w, http.StatusConflict, err.Error())
			case users.ErrPasswordTooLong:
				jsonError(w, http.StatusBadRequest, err.Error())
			default:
				jsonError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		jsonOK(w, u)
	}
}
