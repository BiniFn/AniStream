package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AdminRoutes() {
	h.r.With(middleware.RequireAdmin).Route("/__admin", func(r chi.Router) {
		r.Get("/bulk-job/test", h.testAdminAuth)
		r.Post("/bulk-reprocess-anime-file", h.bulkReprocessAnimeFromFile)
		r.Get("/bulk-job/{jobId}", h.getBulkJobStatus)
		r.Get("/bulk-job/{jobId}/result", h.getBulkJobResult)
		r.Get("/bulk-job/{jobId}/download", h.downloadBulkJobResult)
		r.Get("/bulk-job/{jobId}/failed-ids", h.downloadFailedIds)
		r.Post("/bulk-job/{jobId}/retry", h.retryFailedIds)
		r.Post("/unknown-season-fix", h.unknownSeasonFix)
	})
}

func (h *Handler) testAdminAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) getBulkJobStatus(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	jobID, err := h.pathParam(r, "jobId")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.services.Admin.GetBulkJobStatus(jobID)
	if err != nil {
		log.Error("Failed to get job status", "jobId", jobID, "err", err)
		if err.Error() == "job not found" {
			h.jsonError(w, http.StatusNotFound, "Job not found")
		} else {
			h.jsonError(w, http.StatusInternalServerError, "Failed to get job status")
		}
		return
	}

	h.jsonOK(w, status)
}

func (h *Handler) getBulkJobResult(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	jobID, err := h.pathParam(r, "jobId")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Admin.GetBulkJobResult(jobID)
	if err != nil {
		log.Error("Failed to get job result", "jobId", jobID, "err", err)
		if err.Error() == "job not found" {
			h.jsonError(w, http.StatusNotFound, "Job not found")
		} else {
			h.jsonError(w, http.StatusInternalServerError, "Failed to get job result")
		}
		return
	}

	h.jsonOK(w, result)
}

func (h *Handler) downloadBulkJobResult(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	jobID, err := h.pathParam(r, "jobId")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Admin.GetBulkJobResult(jobID)
	if err != nil {
		log.Error("Failed to get job result", "jobId", jobID, "err", err)
		if err.Error() == "job not found" {
			h.jsonError(w, http.StatusNotFound, "Job not found")
		} else {
			h.jsonError(w, http.StatusInternalServerError, "Failed to get job result")
		}
		return
	}

	filename := "bulk-job-" + jobID + "-result.json"
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		log.Error("Failed to encode result", "err", err)
		return
	}
}

func (h *Handler) downloadFailedIds(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	jobID, err := h.pathParam(r, "jobId")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Admin.GetBulkJobResult(jobID)
	if err != nil {
		log.Error("Failed to get job result", "jobId", jobID, "err", err)
		if err.Error() == "job not found" {
			h.jsonError(w, http.StatusNotFound, "Job not found")
		} else {
			h.jsonError(w, http.StatusInternalServerError, "Failed to get job result")
		}
		return
	}

	filename := "bulk-job-" + jobID + "-failed-ids.json"
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result.FailedIDs); err != nil {
		log.Error("Failed to encode failed IDs", "err", err)
		return
	}
}

func (h *Handler) retryFailedIds(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	jobID, err := h.pathParam(r, "jobId")
	if err != nil {
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Admin.GetBulkJobResult(jobID)
	if err != nil {
		log.Error("Failed to get job result", "jobId", jobID, "err", err)
		if err.Error() == "job not found" {
			h.jsonError(w, http.StatusNotFound, "Job not found")
		} else {
			h.jsonError(w, http.StatusInternalServerError, "Failed to get job result")
		}
		return
	}

	if len(result.FailedIDs) == 0 {
		h.jsonError(w, http.StatusBadRequest, "No failed IDs to retry")
		return
	}

	newJobID, err := h.services.Admin.StartBulkReprocessFromIDsWithParent(result.FailedIDs, jobID)
	if err != nil {
		log.Error("Failed to start retry job", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to start retry job")
		return
	}

	h.jsonOK(w, map[string]any{
		"jobId":   newJobID,
		"message": fmt.Sprintf("Retry job started with %d failed IDs", len(result.FailedIDs)),
	})
}

func (h *Handler) bulkReprocessAnimeFromFile(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Error("Failed to parse multipart form", "err", err)
		h.jsonError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error("Failed to get file from form", "err", err)
		h.jsonError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	if header.Size > 1<<20 {
		h.jsonError(w, http.StatusBadRequest, "File too large, maximum 1MB allowed")
		return
	}

	log.Info("Starting bulk anime reprocessing from file", "filename", header.Filename, "size", header.Size)

	result, err := h.services.Admin.StartBulkReprocessFromFile(r.Context(), file)
	if err != nil {
		log.Error("Failed to start bulk reprocess from file", "err", err)
		h.jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Info("Bulk reprocessing job started from file", "jobId", result.JobID)
	h.jsonOK(w, result)
}

func (h *Handler) unknownSeasonFix(w http.ResponseWriter, r *http.Request) {
	log := h.logger(r)

	count, err := h.deps.Repo.GetAnimeBySeasonCount(r.Context(), repository.SeasonUnknown)
	if err != nil {
		log.Error("Failed to get count of unknown season anime", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to get count of unknown season anime")
		return
	}

	if count == 0 {
		h.jsonOK(w, map[string]any{
			"message": "No anime with unknown season found",
			"count":   0,
		})
		return
	}

	animes, err := h.deps.Repo.GetAnimeBySeason(r.Context(), repository.GetAnimeBySeasonParams{
		Season: repository.SeasonUnknown,
		Limit:  int32(count),
	})
	if err != nil {
		log.Error("Failed to get unknown season anime", "err", err)
		h.jsonError(w, http.StatusInternalServerError, "Failed to get unknown season anime")
		return
	}

	go func() {
		fixedCount := 0
		for _, anime := range animes {
			_, err := h.services.Anime.GetAnimeByID(context.Background(), anime.ID)
			if err != nil {
				log.Error("Failed to reprocess anime for unknown season fix", "animeID", anime.ID, "err", err)
				continue
			}
			fixedCount++
		}
	}()

	h.jsonOK(w, map[string]any{
		"message": fmt.Sprintf("Started reprocessing %d anime with unknown season", count),
		"count":   count,
	})
}
