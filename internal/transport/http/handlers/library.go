package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/library"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) LibraryRoutes() {
	h.r.With(middleware.RequireUser).Route("/library", func(r chi.Router) {
		r.Get("/", h.getLibrary)
		r.Get("/{animeID}", h.getAnimeStatus)
		r.Get("/continue-watching", h.getContinueWatching)
		r.Get("/planning", h.getPlanning)
		r.Post("/{animeID}", h.createLibrary)
		r.Put("/{animeID}", h.updateLibrary)
		r.Delete("/{animeID}", h.deleteAnimeFromLib)

		r.Post("/import", h.importLibrary)
		r.Get("/import/{id}", h.getLibraryImportStatus)
	})
}

// @Summary Get user's anime library
// @Description Get user's anime library
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param status query models.LibraryStatus true "Library status filter"
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.LibraryListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library [get]
func (h *Handler) getLibrary(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		h.jsonError(w, http.StatusBadRequest, "status is required")
		return
	}

	lib, err := h.libraryService.GetLibrary(r.Context(), library.GetLibraryParams{
		UserID:       user.ID,
		Status:       status,
		Page:         page,
		ItemsPerPage: size,
	})

	switch err {
	case library.ErrInvalidStatus:
		h.jsonError(w, http.StatusBadRequest, err.Error())
	case nil:
		h.jsonOK(w, lib)
	default:
		log.Error("failed to get library", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get library")
	}
}

// @Summary Get anime status in library
// @Description Get anime status in library
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param animeID path string true "Anime ID"
// @Success 200 {object} models.LibraryResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/{animeID} [get]
func (h *Handler) getAnimeStatus(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	animeID, err := h.pathParam(r, "animeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.libraryService.GetLibraryByAnimeID(r.Context(), user.ID, animeID)
	switch err {
	case library.ErrLibraryNotFound:
		h.jsonError(w, http.StatusNotFound, "library not found")
	case nil:
		h.jsonOK(w, status)
	default:
		log.Error("failed to get library by anime ID", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get library by anime ID")
		return
	}
}

// @Summary Get continue watching list
// @Description Get continue watching list
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.LibraryListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/continue-watching [get]
func (h *Handler) getContinueWatching(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	library, err := h.libraryService.GetContinueWatching(r.Context(), library.GetContinueWatchingAnimeParams{
		UserID:       user.ID,
		Page:         page,
		ItemsPerPage: size,
	})
	if err != nil {
		log.Error("failed to get continue watching", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get continue watching")
		return
	}

	h.jsonOK(w, library)
}

// @Summary Get plan to watch list
// @Description Get plan to watch list
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param page query int false "Page number"
// @Param itemsPerPage query int false "Number of items per page"
// @Success 200 {object} models.LibraryListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/planning [get]
func (h *Handler) getPlanning(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	library, err := h.libraryService.GetPlanToWatch(r.Context(), library.GetPlanToWatchAnimeParams{
		UserID:       user.ID,
		Page:         page,
		ItemsPerPage: size,
	})
	if err != nil {
		log.Error("failed to get plan to watch", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get plan to watch")
		return
	}

	h.jsonOK(w, library)
}

// @Summary Remove anime from library
// @Description Remove anime from library
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param animeID path string true "Anime ID"
// @Success 200
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/{animeID} [delete]
func (h *Handler) deleteAnimeFromLib(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	animeID, err := h.pathParam(r, "animeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.libraryService.DeleteLibrary(r.Context(), user.ID, animeID)
	if err != nil {
		log.Error("failed to delete anime from library", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to delete anime from library")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Add anime to library
// @Description Add anime to library
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param animeID path string true "Anime ID"
// @Param library body models.LibraryRequest true "Library object"
// @Success 200 {object} models.LibraryResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/{animeID} [post]
func (h *Handler) createLibrary(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	animeID, err := h.pathParam(r, "animeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.LibraryRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	lib, err := h.libraryService.CreateLibrary(r.Context(), user.ID, animeID, string(req.Status), req.WatchedEpisodes)
	switch err {
	case library.ErrInvalidStatus, library.ErrInvalidWatchedEpisodes:
		h.jsonError(w, http.StatusBadRequest, err.Error())
	case nil:
		h.jsonOK(w, lib)
	default:
		log.Error("failed to save anime to library", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to save anime to library")
	}
}

// @Summary Update anime in library
// @Description Update anime in library
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param animeID path string true "Anime ID"
// @Param library body models.LibraryRequest true "Library object"
// @Success 200 {object} models.LibraryResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/{animeID} [put]
func (h *Handler) updateLibrary(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	animeID, err := h.pathParam(r, "animeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.LibraryRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	lib, err := h.libraryService.UpdateLibrary(r.Context(), user.ID, animeID, string(req.Status), req.WatchedEpisodes)
	switch err {
	case library.ErrInvalidStatus, library.ErrInvalidWatchedEpisodes:
		h.jsonError(w, http.StatusBadRequest, err.Error())
	case nil:
		h.jsonOK(w, lib)
	default:
		log.Error("failed to update anime in library", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to update anime in library")
	}
}

// @Summary Import library from external provider
// @Description Import library from external provider
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param provider query string true "External provider to import from"
// @Success 200 {object} models.ImportJobResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/import [post]
func (h *Handler) importLibrary(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		h.jsonError(w, http.StatusBadRequest, "provider is required")
		return
	}

	id, err := h.libraryService.ImportLibrary(r.Context(), user.ID, provider)
	switch err {
	case library.ErrInvalidProvider:
		h.jsonError(w, http.StatusBadRequest, err.Error())
	case nil:
		h.jsonOK(w, models.ImportJobResponse{ID: id})
	default:
		log.Error("failed to import library", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to import library")
	}
}

// @Summary Get library import status
// @Description Get library import status
// @Tags Library
// @Accept json
// @Produce json
// @Security cookieAuth
// @Param id path string true "Import job ID"
// @Success 200 {object} models.LibraryImportJobResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /library/import/{id} [get]
func (h *Handler) getLibraryImportStatus(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	id, err := h.pathParam(r, "id")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.libraryService.GetImportLibraryStatus(r.Context(), id)
	switch err {
	case library.ErrJobNotFound:
		h.jsonError(w, http.StatusNotFound, err.Error())
	case nil:
		h.jsonOK(w, status)
	default:
		log.Error("failed to get library import status", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get library import status")
	}
}
