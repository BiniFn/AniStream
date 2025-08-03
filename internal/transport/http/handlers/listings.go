package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AnimeListingRoutes() {
	h.r.Route("/anime", func(r chi.Router) {
		r.Get("/recently-updated", h.listRecentlyUpdated)
		r.Get("/seasonal", h.seasonalAnimes)
		r.Get("/random", h.randomAnime)
		r.Get("/genres", h.listGenres)
		r.Get("/genres/{genre}", h.animeByGenre)
		r.Get("/search", h.searchAnimes)
		r.Get("/trending", h.trendingAnimes)
		r.Get("/popular", h.popularAnimes)
	})
}

func (h *Handler) listRecentlyUpdated(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetRecentlyUpdatedAnimes(r.Context(), page, size)
	if err != nil {
		log.Error("failed to fetch recently updated animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch recently updated animes")
		return
	}
	h.jsonOK(w, resp)
}

func (h *Handler) seasonalAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	resp, err := h.animeService.GetSeasonalAnimes(r.Context())
	if err != nil {
		log.Error("failed to fetch seasonal animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch seasonal animes")
		return
	}
	h.jsonOK(w, resp)
}

func (h *Handler) randomAnime(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	genre := r.URL.Query().Get("genre")

	var (
		resp anime.AnimeDto
		err  error
	)

	if genre != "" {
		resp, err = h.animeService.GetRandomAnimeByGenre(r.Context(), genre)
		if err != nil {
			log.Error("failed to fetch random anime by genre", "genre", genre, "err", err)
			h.jsonError(w, http.StatusInternalServerError, "failed to fetch random anime by genre")
			return
		}
	} else {
		resp, err = h.animeService.GetRandomAnime(r.Context())
		if err != nil {
			log.Error("failed to fetch random animes", "err", err)
			h.jsonError(w, http.StatusInternalServerError, "failed to fetch random animes")
			return
		}
	}

	h.jsonOK(w, resp)
}

func (h *Handler) listGenres(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)
	resp, err := h.animeService.GetAnimeGenres(r.Context())
	if err != nil {
		log.Error("failed to fetch anime genres", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime genres")
		return
	}
	h.jsonOK(w, resp)
}

func (h *Handler) searchAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	query := r.URL.Query().Get("q")
	if query == "" {
		h.jsonError(w, http.StatusBadRequest, "search query is required")
		return
	}
	if len(query) < 3 {
		h.jsonError(w, http.StatusBadRequest, "search query must be at least 3 characters long")
		return
	}

	genre := r.URL.Query().Get("genre")
	if genre != "" && len(genre) < 3 {
		h.jsonError(w, http.StatusBadRequest, "genre query must be at least 3 characters long")
		return
	}

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.SearchAnimes(r.Context(), query, genre, page, size)
	if err != nil {
		log.Error("failed to search animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to search animes")
		return
	}
	h.jsonOK(w, resp)
}

func (h *Handler) animeByGenre(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	genre, err := h.pathParam(r, "genre")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	page, size, err := h.parsePagination(r, 1, 30)
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.animeService.GetAnimesByGenre(r.Context(), genre, page, size)
	if err != nil {
		log.Error("failed to fetch animes by genre", "genre", genre, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch animes by genre")
		return
	}
	h.jsonOK(w, resp)
}

func (h *Handler) trendingAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	resp, err := h.animeService.GetTrendingAnimes(r.Context())
	if err != nil {
		log.Error("failed to fetch trending animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch trending animes")
		return
	}
	h.jsonOK(w, resp)
}

func (h *Handler) popularAnimes(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	resp, err := h.animeService.GetPopularAnimes(r.Context())
	if err != nil {
		log.Error("failed to fetch popular animes", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch popular animes")
		return
	}
	h.jsonOK(w, resp)
}
