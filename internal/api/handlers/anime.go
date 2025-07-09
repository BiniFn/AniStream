package handlers

import (
	"net/http"

	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeRoutes(r chi.Router, svc *animeSvc.Service) {
	r.Get("/recently-updated", listRecentlyUpdated(svc))
}

func listRecentlyUpdated(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetRecentlyUpdatedAnimes(r.Context(), page, size)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to fetch recently updated animes")
			return
		}
		jsonOK(w, resp)
	}
}
