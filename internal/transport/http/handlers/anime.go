package handlers

import (
	"net/http"
	"sync"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/anime"
	"github.com/coeeter/aniways/internal/service/library"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type GetAnimeByIDInput struct {
	ID string `in:"path=id"`
}

func (h *Handler) AnimeDetailsRoutes() {
	h.r.With(httpin.NewInput(GetAnimeByIDInput{})).Route("/anime/{id}", func(r chi.Router) {
		r.Get("/", h.getAnimeByID)
		r.Get("/full", h.getAnimeFull)
		r.Get("/variations", h.getAnimeVariations)
		r.Get("/trailer", h.getAnimeTrailer)
		r.Get("/banner", h.getAnimeBanner)
		r.Get("/franchise", h.getAnimeFranchise)
		r.Get("/characters", h.getAnimeCharacters)
	})
}

// @Summary Get anime by ID
// @Description Get anime by ID
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.AnimeWithMetadataResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id} [get]
func (h *Handler) getAnimeByID(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.services.Anime.GetAnimeByID(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

// @Summary Get anime trailer
// @Description Get anime trailer
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.TrailerResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/trailer [get]
func (h *Handler) getAnimeTrailer(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.services.Anime.GetAnimeTrailer(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case anime.TrailerNotFound:
		log.Warn("trailer not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "trailer not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

// @Summary Get anime banner
// @Description Get anime banner
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.BannerResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/banner [get]
func (h *Handler) getAnimeBanner(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.services.Anime.GetAnimeBanner(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case anime.BannerNotFound:
		log.Warn("banner not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "banner not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

// @Summary Get anime franchise relations
// @Description Get anime franchise relations
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.RelationsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/franchise [get]
func (h *Handler) getAnimeFranchise(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.services.Anime.GetAnimeRelations(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

// @Summary Get anime characters
// @Description Get anime characters
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.CharactersResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/characters [get]
func (h *Handler) getAnimeCharacters(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	resp, err := h.services.Anime.GetAnimeCharacters(r.Context(), id)
	switch err {
	case anime.ErrAnimeNotFound:
		log.Warn("anime not found", "id", id, "err", err)
		h.jsonError(w, http.StatusNotFound, "anime not found")
		return
	case nil:
		h.jsonOK(w, resp)
		return
	default:
		log.Error("failed to fetch anime details", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
		return
	}
}

// @Summary Get anime variations
// @Description Get all variations of an anime (same MAL ID)
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {array} models.AnimeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/variations [get]
func (h *Handler) getAnimeVariations(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	variations, err := h.services.Anime.GetAnimeVariations(r.Context(), id)
	if err != nil {
		log.Error("failed to fetch anime variations", "id", id, "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime variations")
		return
	}

	h.jsonOK(w, variations)
}

// @Summary Get full anime details
// @Description Get all anime details in a single response including anime data, banner, trailer, episodes, franchise, characters, and library status (if authenticated)
// @Tags Anime
// @Accept json
// @Produce json
// @Param id path string true "Anime ID"
// @Success 200 {object} models.AnimeFullResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /anime/{id}/full [get]
func (h *Handler) getAnimeFull(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	input := h.getHttpInput(r).(*GetAnimeByIDInput)
	id := input.ID

	var (
		animeData     *models.AnimeWithMetadataResponse
		banner        *models.BannerResponse
		trailer       *models.TrailerResponse
		episodes      models.EpisodeListResponse
		franchise     *models.RelationsResponse
		characters    models.CharactersResponse
		variations    []models.AnimeResponse
		libStatus     *models.LibraryResponse
		animeErr      error
		bannerErr     error
		trailerErr    error
		episodesErr   error
		franchiseErr  error
		charactersErr error
		variationsErr error
		libErr        error
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		animeData, animeErr = h.services.Anime.GetAnimeByID(r.Context(), id)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bannerResp, err := h.services.Anime.GetAnimeBanner(r.Context(), id)
		if err == nil {
			banner = &bannerResp
		} else if err != anime.BannerNotFound {
			bannerErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		trailerResp, err := h.services.Anime.GetAnimeTrailer(r.Context(), id)
		if err == nil {
			trailer = trailerResp
		} else if err != anime.TrailerNotFound {
			trailerErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		episodes, episodesErr = h.services.Anime.GetAnimeEpisodes(r.Context(), id)
		if episodesErr != nil {
			episodes = nil
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		franchiseResp, err := h.services.Anime.GetAnimeRelations(r.Context(), id)
		if err == nil {
			franchise = &franchiseResp
		} else {
			franchiseErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		characters, charactersErr = h.services.Anime.GetAnimeCharacters(r.Context(), id)
		if charactersErr != nil {
			characters = nil
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		variations, variationsErr = h.services.Anime.GetAnimeVariations(r.Context(), id)
		if variationsErr != nil {
			variations = nil
		}
	}()

	user := middleware.GetUser(r)
	if user != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			libResp, err := h.services.Library.GetLibraryByAnimeID(r.Context(), user.ID, id)
			if err == nil {
				libStatus = &libResp
			} else if err != library.ErrLibraryNotFound {
				libErr = err
			}
		}()
	}

	wg.Wait()

	if animeErr != nil {
		switch animeErr {
		case anime.ErrAnimeNotFound:
			log.Warn("anime not found", "id", id, "err", animeErr)
			h.jsonError(w, http.StatusNotFound, "anime not found")
			return
		default:
			log.Error("failed to fetch anime details", "id", id, "err", animeErr)
			h.jsonError(w, http.StatusInternalServerError, "failed to fetch anime details")
			return
		}
	}

	if bannerErr != nil {
		log.Warn("failed to fetch banner", "id", id, "err", bannerErr)
	}
	if trailerErr != nil && trailerErr != anime.TrailerNotFound {
		log.Warn("failed to fetch trailer", "id", id, "err", trailerErr)
	}
	if episodesErr != nil {
		log.Warn("failed to fetch episodes", "id", id, "err", episodesErr)
	}
	if franchiseErr != nil {
		log.Warn("failed to fetch franchise", "id", id, "err", franchiseErr)
	}
	if charactersErr != nil {
		log.Warn("failed to fetch characters", "id", id, "err", charactersErr)
	}
	if variationsErr != nil {
		log.Warn("failed to fetch variations", "id", id, "err", variationsErr)
	}
	if libErr != nil {
		log.Warn("failed to fetch library status", "id", id, "err", libErr)
	}

	response := models.AnimeFullResponse{
		Anime:         *animeData,
		Banner:        banner,
		Trailer:       trailer,
		Episodes:      episodes,
		Franchise:     franchise,
		LibraryStatus: libStatus,
		Characters:    characters,
		Variations:    variations,
	}

	h.jsonOK(w, response)
}
