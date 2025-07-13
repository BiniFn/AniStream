package handlers

import (
	"log"
	"net/http"

	"github.com/coeeter/aniways/internal/models"
	animeSvc "github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeListingsRoutes(r chi.Router, svc *animeSvc.Service) {
	r.Get("/recently-updated", listRecentlyUpdated(svc))
	r.Get("/seasonal", seasonalAnimes(svc))
	r.Get("/random", randomAnime(svc))
	r.Get("/genres", listGenres(svc))
	r.Get("/genres/{genre}", animeByGenre(svc))
	r.Get("/search", searchAnimes(svc))
	r.Get("/trending", trendingAnimes(svc))
	r.Get("/popular", popularAnimes(svc))
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
			log.Printf("failed to fetch recently updated animes: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch recently updated animes")
			return
		}
		jsonOK(w, resp)
	}
}

func seasonalAnimes(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.GetSeasonalAnimes(r.Context())
		if err != nil {
			log.Printf("failed to fetch seasonal animes: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch seasonal animes")
			return
		}
		jsonOK(w, resp)
	}
}

func randomAnime(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genre := r.URL.Query().Get("genre")

		var (
			resp models.AnimeDto
			err  error
		)

		if genre != "" {
			resp, err = svc.GetRandomAnimeByGenre(r.Context(), genre)
			if err != nil {
				log.Printf("failed to fetch random anime by genre %s: %v", genre, err)
				jsonError(w, http.StatusInternalServerError, "failed to fetch random anime by genre")
				return
			}
		} else {
			resp, err = svc.GetRandomAnime(r.Context())
			if err != nil {
				log.Printf("failed to fetch random animes: %v", err)
				jsonError(w, http.StatusInternalServerError, "failed to fetch random animes")
				return
			}
		}

		jsonOK(w, resp)
	}
}

func listGenres(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.GetAnimeGenres(r.Context())
		if err != nil {
			log.Printf("failed to fetch anime genres: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime genres")
			return
		}
		jsonOK(w, resp)
	}
}

func searchAnimes(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			jsonError(w, http.StatusBadRequest, "search query is required")
			return
		}
		if len(query) < 3 {
			jsonError(w, http.StatusBadRequest, "search query must be at least 3 characters long")
			return
		}

		genre := r.URL.Query().Get("genre")
		if genre != "" && len(genre) < 3 {
			jsonError(w, http.StatusBadRequest, "genre query must be at least 3 characters long")
			return
		}

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.SearchAnimes(r.Context(), query, genre, page, size)
		if err != nil {
			log.Printf("failed to search animes: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to search animes")
			return
		}
		jsonOK(w, resp)
	}
}

func animeByGenre(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		genre := chi.URLParam(r, "genre")
		if genre == "" {
			jsonError(w, http.StatusBadRequest, "genre is required")
			return
		}

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimesByGenre(r.Context(), genre, page, size)
		if err != nil {
			log.Printf("failed to fetch animes by genre %s: %v", genre, err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch animes by genre")
			return
		}
		jsonOK(w, resp)
	}
}

func trendingAnimes(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.GetTrendingAnimes(r.Context())
		if err != nil {
			log.Printf("failed to fetch trending animes: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch trending animes")
			return
		}
		jsonOK(w, resp)
	}
}

func popularAnimes(svc *animeSvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.GetPopularAnimes(r.Context())
		if err != nil {
			log.Printf("failed to fetch popular animes: %v", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch popular animes")
			return
		}
		jsonOK(w, resp)
	}
}
