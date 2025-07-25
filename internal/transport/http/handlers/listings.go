package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func MountAnimeListingsRoutes(r chi.Router, svc *anime.AnimeService) {
	r.Get("/recently-updated", listRecentlyUpdated(svc))
	r.Get("/seasonal", seasonalAnimes(svc))
	r.Get("/random", randomAnime(svc))
	r.Get("/genres", listGenres(svc))
	r.Get("/genres/{genre}", animeByGenre(svc))
	r.Get("/search", searchAnimes(svc))
	r.Get("/trending", trendingAnimes(svc))
	r.Get("/popular", popularAnimes(svc))
}

func listRecentlyUpdated(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetRecentlyUpdatedAnimes(r.Context(), page, size)
		if err != nil {
			log.Error("failed to fetch recently updated animes", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch recently updated animes")
			return
		}
		jsonOK(w, resp)
	}
}

func seasonalAnimes(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		resp, err := svc.GetSeasonalAnimes(r.Context())
		if err != nil {
			log.Error("failed to fetch seasonal animes", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch seasonal animes")
			return
		}
		jsonOK(w, resp)
	}
}

func randomAnime(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		genre := r.URL.Query().Get("genre")

		var (
			resp anime.AnimeDto
			err  error
		)

		if genre != "" {
			resp, err = svc.GetRandomAnimeByGenre(r.Context(), genre)
			if err != nil {
				log.Error("failed to fetch random anime by genre", "genre", genre, "err", err)
				jsonError(w, http.StatusInternalServerError, "failed to fetch random anime by genre")
				return
			}
		} else {
			resp, err = svc.GetRandomAnime(r.Context())
			if err != nil {
				log.Error("failed to fetch random animes", "err", err)
				jsonError(w, http.StatusInternalServerError, "failed to fetch random animes")
				return
			}
		}

		jsonOK(w, resp)
	}
}

func listGenres(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)
		resp, err := svc.GetAnimeGenres(r.Context())
		if err != nil {
			log.Error("failed to fetch anime genres", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch anime genres")
			return
		}
		jsonOK(w, resp)
	}
}

func searchAnimes(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

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
			log.Error("failed to search animes", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to search animes")
			return
		}
		jsonOK(w, resp)
	}
}

func animeByGenre(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		genre, err := pathParam(r, "genre")
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		page, size, err := parsePagination(r, 1, 30)
		if err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := svc.GetAnimesByGenre(r.Context(), genre, page, size)
		if err != nil {
			log.Error("failed to fetch animes by genre", "genre", genre, "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch animes by genre")
			return
		}
		jsonOK(w, resp)
	}
}

func trendingAnimes(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		resp, err := svc.GetTrendingAnimes(r.Context())
		if err != nil {
			log.Error("failed to fetch trending animes", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch trending animes")
			return
		}
		jsonOK(w, resp)
	}
}

func popularAnimes(svc *anime.AnimeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger(r)

		resp, err := svc.GetPopularAnimes(r.Context())
		if err != nil {
			log.Error("failed to fetch popular animes", "err", err)
			jsonError(w, http.StatusInternalServerError, "failed to fetch popular animes")
			return
		}
		jsonOK(w, resp)
	}
}
