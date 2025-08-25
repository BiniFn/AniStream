package handlers

import (
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/users"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) UserRoutes() {
	h.r.Route("/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.With(middleware.RequireUser).Group(func(r chi.Router) {
			r.Put("/", h.updateUser)
			r.Delete("/", h.deleteUser)
			r.Put("/password", h.updatePassword)
			r.Put("/image", h.updateImage)
		})
	})
}

// @Summary Create new user
// @Description Create new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User object"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [post]
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	var req models.CreateUserRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	u, err := h.userService.CreateUser(r.Context(), req.Username, req.Email, req.Password)
	switch err {
	case nil:
		h.jsonOK(w, u)
	case users.ErrEmailTaken, users.ErrUsernameTaken:
		log.Warn("username or email already taken", "username", req.Username, "email", req.Email)
		h.jsonError(w, http.StatusConflict, err.Error())
	case users.ErrPasswordTooLong:
		log.Warn("password too long", "username", req.Username)
		h.jsonError(w, http.StatusBadRequest, err.Error())
	default:
		log.Error("user creation failed", "err", err)
		h.jsonError(w, http.StatusInternalServerError, err.Error())
	}
}

// @Summary Update user password
// @Description Update user password
// @Tags Users
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param password body models.UpdatePasswordRequest true "Password object"
// @Success 200
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/password [put]
func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	var req models.UpdatePasswordRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	err := h.userService.UpdatePassword(r.Context(), user.ID, req.OldPassword, req.NewPassword)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case users.ErrPasswordTooLong:
		log.Warn("password too long", "id", user.ID)
		h.jsonError(w, http.StatusBadRequest, err.Error())
	case users.ErrInvalidAuth:
		log.Warn("invalid password", "id", user.ID)
		h.jsonError(w, http.StatusUnauthorized, err.Error())
	default:
		log.Error("password update failed", "err", err)
		h.jsonError(w, http.StatusInternalServerError, err.Error())
	}
}

// @Summary Update user information
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param user body models.UpdateUserRequest true "User object"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [put]
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	var req models.UpdateUserRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	u, err := h.userService.UpdateUser(r.Context(), user.ID, req.Username, req.Email)
	switch err {
	case nil:
		h.jsonOK(w, u)
	case users.ErrEmailTaken, users.ErrUsernameTaken:
		log.Warn("username or email already taken", "username", req.Username, "email", req.Email)
		h.jsonError(w, http.StatusConflict, err.Error())
	default:
		log.Error("user update failed", "err", err)
		h.jsonError(w, http.StatusInternalServerError, err.Error())
	}
}

// @Summary Delete user account
// @Description Delete user account
// @Tags Users
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param password body models.DeleteUserRequest true "Password object"
// @Success 204
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [delete]
func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	var req models.DeleteUserRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	err := h.userService.DeleteUser(r.Context(), user.ID, req.Password)
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
		h.jsonError(w, http.StatusUnauthorized, err.Error())
	default:
		log.Error("user deletion failed", "err", err)
		h.jsonError(w, http.StatusInternalServerError, err.Error())
	}

}

// @Summary Update user profile picture
// @Description Update user profile picture
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Security cookieAuth
// @Param image formData file true "Image file"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/image [put]
func (h *Handler) updateImage(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	if err := r.ParseMultipartForm(10 << 20 /* 10 MiB */); err != nil {
		h.jsonError(w, http.StatusBadRequest, "bad request")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, "image field is required")
		return
	}
	defer file.Close()

	err = h.userService.UpdateProfilePicture(r.Context(), user.ID, file)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		log.Error("user image update failed", "id", user.ID, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "internal error")
	}
}
