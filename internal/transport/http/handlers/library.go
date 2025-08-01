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
		if err != nil {
			log.Error("failed to get library by anime ID", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to get library by anime ID")
			return
		}

		jsonOK(w, status)
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

		jsonOK(w, nil)
	}
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

		var body struct {
			Status         string `json:"status"`
			WatchedEpisode int32  `json:"watchedEpisodes"`
		}
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

		var body struct {
			Status         string `json:"status"`
			WatchedEpisode int32  `json:"watched_episode"`
		}
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
