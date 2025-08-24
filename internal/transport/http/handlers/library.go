package handlers

import (
	"encoding/json"
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


func (h *Handler) createLibrary(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	animeID, err := h.pathParam(r, "animeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.LibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	lib, err := h.libraryService.CreateLibrary(r.Context(), user.ID, animeID, req.Status, req.WatchedEpisodes)
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

func (h *Handler) updateLibrary(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	user := middleware.GetUser(r)

	animeID, err := h.pathParam(r, "animeID")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.LibraryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	lib, err := h.libraryService.UpdateLibrary(r.Context(), user.ID, animeID, req.Status, req.WatchedEpisodes)
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
