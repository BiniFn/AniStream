package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coeeter/aniways/internal/service/library"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func MountLibraryRoutes(r chi.Router, svc *library.LibraryService) {
	r.With(middleware.RequireUser).Group(func(r chi.Router) {
		r.Get("/", getLibrary(svc))
		r.Get("/{animeID}", getAnimeStatus(svc))
		r.Get("/continue-watching", getContinueWatching(svc))
		r.Get("/planning", getPlanning(svc))
		r.Post("/{animeID}", createLibrary(svc))
		r.Put("/{animeID}", updateLibrary(svc))
		r.Delete("/{animeID}", deleteAnimeFromLib(svc))

		r.Post("/import", importLibrary(svc))
		r.Get("/import/{id}", getLibraryImportStatus(svc))
	})
}

func getLibrary(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		status := r.URL.Query().Get("status")
		if status == "" {
			jsonError(w, http.StatusBadRequest, "status is required")
			return
		}

		lib, err := svc.GetLibrary(r.Context(), library.GetLibraryParams{
			UserID:       user.ID,
			Status:       status,
			Page:         page,
			ItemsPerPage: size,
		})

		switch err {
		case library.ErrInvalidStatus:
			jsonError(w, http.StatusBadRequest, err.Error())
		case nil:
			jsonOK(w, lib)
		default:
			log.Error("failed to get library", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get library")
		}

	}
}

func getAnimeStatus(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		animeID, err := pathParam(r, "animeID")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		status, err := svc.GetLibraryByAnimeID(r.Context(), user.ID, animeID)
		switch err {
		case library.ErrLibraryNotFound:
			jsonError(w, http.StatusNotFound, "library not found")
		case nil:
			jsonOK(w, status)
		default:
			log.Error("failed to get library by anime ID", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get library by anime ID")
			return
		}
	}
}

func getContinueWatching(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		library, err := svc.GetContinueWatching(r.Context(), library.GetContinueWatchingAnimeParams{
			UserID:       user.ID,
			Page:         page,
			ItemsPerPage: size,
		})
		if err != nil {
			log.Error("failed to get continue watching", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get continue watching")
			return
		}

		jsonOK(w, library)
	}
}

func getPlanning(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		library, err := svc.GetPlanToWatch(r.Context(), library.GetPlanToWatchAnimeParams{
			UserID:       user.ID,
			Page:         page,
			ItemsPerPage: size,
		})
		if err != nil {
			log.Error("failed to get plan to watch", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get plan to watch")
			return
		}

		jsonOK(w, library)
	}
}

func deleteAnimeFromLib(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		animeID, err := pathParam(r, "animeID")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = svc.DeleteLibrary(r.Context(), user.ID, animeID)
		if err != nil {
			log.Error("failed to delete anime from library", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to delete anime from library")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type LibraryBody struct {
	Status         string `json:"status"`
	WatchedEpisode int32  `json:"watchedEpisodes"`
}

func createLibrary(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		animeID, err := pathParam(r, "animeID")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		var body LibraryBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		lib, err := svc.CreateLibrary(r.Context(), user.ID, animeID, body.Status, body.WatchedEpisode)
		switch err {
		case library.ErrInvalidStatus, library.ErrInvalidWatchedEpisodes:
			jsonError(w, http.StatusBadRequest, err.Error())
		case nil:
			jsonOK(w, lib)
		default:
			log.Error("failed to save anime to library", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to save anime to library")
		}
	}
}

func updateLibrary(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		animeID, err := pathParam(r, "animeID")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		var body LibraryBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		lib, err := svc.UpdateLibrary(r.Context(), user.ID, animeID, body.Status, body.WatchedEpisode)
		switch err {
		case library.ErrInvalidStatus, library.ErrInvalidWatchedEpisodes:
			jsonError(w, http.StatusBadRequest, err.Error())
		case nil:
			jsonOK(w, lib)
		default:
			log.Error("failed to update anime in library", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to update anime in library")
		}
	}
}

func importLibrary(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		user := middleware.GetUser(r)

		provider := r.URL.Query().Get("provider")
		if provider == "" {
			jsonError(w, http.StatusBadRequest, "provider is required")
			return
		}

		id, err := svc.ImportLibrary(r.Context(), user.ID, provider)
		switch err {
		case library.ErrInvalidProvider:
			jsonError(w, http.StatusBadRequest, err.Error())
		case nil:
			jsonOK(w, map[string]string{
				"id": id,
			})
		default:
			log.Error("failed to import library", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to import library")
		}
	}
}

func getLibraryImportStatus(svc *library.LibraryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		id, err := pathParam(r, "id")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		status, err := svc.GetImportLibraryStatus(r.Context(), id)
		switch err {
		case library.ErrJobNotFound:
			jsonError(w, http.StatusNotFound, err.Error())
		case nil:
			jsonOK(w, status)
		default:
			log.Error("failed to get library import status", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get library import status")
		}
	}
}
