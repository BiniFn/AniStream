package handlers

import (
	"net/http"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/service/desktop"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) DesktopRoutes() {
	h.r.Route("/desktop/releases", func(r chi.Router) {
		r.Get("/", h.getAllDesktopReleases)
		r.Get("/latest", h.getLatestDesktopRelease)
		r.Get("/{version}", h.getDesktopReleaseByVersion)

		r.With(middleware.RequireDesktopReleaseKey(h.deps.Env)).Post("/", h.createDesktopRelease)
		r.With(middleware.RequireDesktopReleaseKey(h.deps.Env)).Delete("/{version}", h.deleteDesktopRelease)
	})
}

// @Summary Get all desktop releases
// @Description Get all desktop releases grouped by version
// @Tags Desktop
// @Accept json
// @Produce json
// @Success 200 {array} models.DesktopReleaseResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /desktop/releases [get]
func (h *Handler) getAllDesktopReleases(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	releases, err := h.services.Desktop.GetAllReleases(r.Context())
	if err != nil {
		log.Error("failed to get desktop releases", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get desktop releases")
		return
	}

	h.jsonOK(w, releases)
}

// @Summary Get latest desktop release
// @Description Get the latest desktop release with all platform binaries
// @Tags Desktop
// @Accept json
// @Produce json
// @Success 200 {object} models.DesktopVersionResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /desktop/releases/latest [get]
func (h *Handler) getLatestDesktopRelease(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	release, err := h.services.Desktop.GetLatestReleases(r.Context())
	if err != nil {
		log.Error("failed to get latest desktop release", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get latest desktop release")
		return
	}

	if release == nil {
		h.jsonError(w, http.StatusNotFound, "no releases found")
		return
	}

	h.jsonOK(w, release)
}

// @Summary Get desktop release by version
// @Description Get a specific desktop release version with all platform binaries
// @Tags Desktop
// @Accept json
// @Produce json
// @Param version path string true "Version string (e.g. 1.0.0)"
// @Success 200 {object} models.DesktopVersionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /desktop/releases/{version} [get]
func (h *Handler) getDesktopReleaseByVersion(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	version, err := h.pathParam(r, "version")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	release, err := h.services.Desktop.GetReleasesByVersion(r.Context(), version)
	if err != nil {
		log.Error("failed to get desktop release", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to get desktop release")
		return
	}

	if release == nil {
		h.jsonError(w, http.StatusNotFound, "release not found")
		return
	}

	h.jsonOK(w, release)
}

// @Summary Create desktop release
// @Description Create a new desktop release for a specific platform
// @Tags Desktop
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param release body models.CreateDesktopReleaseRequest true "Release object"
// @Success 201 {object} models.DesktopReleaseResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /desktop/releases [post]
func (h *Handler) createDesktopRelease(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	var req models.CreateDesktopReleaseRequest
	if !h.parseAndValidate(w, r, &req) {
		return
	}

	release, err := h.services.Desktop.CreateRelease(r.Context(), desktop.CreateReleaseParams{
		Version:      req.Version,
		Platform:     req.Platform,
		DownloadURL:  req.DownloadURL,
		FileName:     req.FileName,
		FileSize:     req.FileSize,
		ReleaseNotes: req.ReleaseNotes,
	})
	if err != nil {
		log.Error("failed to create desktop release", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to create desktop release")
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.jsonOK(w, release)
}

// @Summary Delete desktop releases by version
// @Description Delete all desktop releases for a specific version
// @Tags Desktop
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param version path string true "Version string (e.g. 1.0.0)"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /desktop/releases/{version} [delete]
func (h *Handler) deleteDesktopRelease(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	version, err := h.pathParam(r, "version")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Desktop.DeleteReleasesByVersion(r.Context(), version)
	if err != nil {
		log.Error("failed to delete desktop release", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "failed to delete desktop release")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
