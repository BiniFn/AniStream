package models

import (
	"sync"
	"time"
)

type ReprocessResult struct {
	HiAnimeID string `json:"hiAnimeId" example:"anime-id-1"`
	Success   bool   `json:"success" example:"true"`
	Message   string `json:"message" example:"Successfully updated anime details"`
}

type ReprocessError struct {
	HiAnimeID string `json:"hiAnimeId" example:"anime-id-2"`
	Error     string `json:"error" example:"Failed to fetch anime details"`
}

type BulkReprocessResponse struct {
	JobID   string `json:"jobId" example:"job-12345"`
	Message string `json:"message" example:"Bulk reprocessing job started"`
}

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type BulkJob struct {
	ID         string
	Status     JobStatus
	Total      int
	Processed  int
	Success    int
	Failed     int
	Message    string
	StartedAt  time.Time
	UpdatedAt  time.Time
	HiAnimeIDs []string
	Results    []ReprocessResult
	Errors     []ReprocessError
	Mu         sync.RWMutex
}

type BulkJobStatus struct {
	JobID     string    `json:"jobId" example:"job-12345"`
	Status    JobStatus `json:"status" example:"running"`
	Total     int       `json:"total" example:"1000"`
	Processed int       `json:"processed" example:"250"`
	Success   int       `json:"success" example:"200"`
	Failed    int       `json:"failed" example:"50"`
	Progress  float64   `json:"progress" example:"25.0"`
	Message   string    `json:"message" example:"Processing anime details..."`
	StartedAt string    `json:"startedAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string    `json:"updatedAt" example:"2023-01-01T00:05:00Z"`
}

type BulkJobResult struct {
	JobID     string             `json:"jobId"`
	Status    JobStatus          `json:"status"`
	Total     int                `json:"total"`
	Processed int                `json:"processed"`
	Success   int                `json:"success"`
	Failed    int                `json:"failed"`
	StartedAt string             `json:"startedAt"`
	UpdatedAt string             `json:"updatedAt"`
	Results   []ReprocessResult  `json:"results"`
	Errors    []ReprocessError   `json:"errors"`
	FailedIDs []string           `json:"failedIds"`
}

func (job *BulkJob) GetProgress() float64 {
	job.Mu.RLock()
	defer job.Mu.RUnlock()

	if job.Total == 0 {
		return 0
	}
	return float64(job.Processed) / float64(job.Total) * 100
}

func (job *BulkJob) ToStatusResponse() BulkJobStatus {
	job.Mu.RLock()
	defer job.Mu.RUnlock()

	return BulkJobStatus{
		JobID:     job.ID,
		Status:    job.Status,
		Total:     job.Total,
		Processed: job.Processed,
		Success:   job.Success,
		Failed:    job.Failed,
		Progress:  job.GetProgress(),
		Message:   job.Message,
		StartedAt: job.StartedAt.Format(time.RFC3339),
		UpdatedAt: job.UpdatedAt.Format(time.RFC3339),
	}
}

func (job *BulkJob) ToResultResponse() BulkJobResult {
	job.Mu.RLock()
	defer job.Mu.RUnlock()

	var failedIDs []string
	for _, err := range job.Errors {
		failedIDs = append(failedIDs, err.HiAnimeID)
	}

	return BulkJobResult{
		JobID:     job.ID,
		Status:    job.Status,
		Total:     job.Total,
		Processed: job.Processed,
		Success:   job.Success,
		Failed:    job.Failed,
		StartedAt: job.StartedAt.Format(time.RFC3339),
		UpdatedAt: job.UpdatedAt.Format(time.RFC3339),
		Results:   job.Results,
		Errors:    job.Errors,
		FailedIDs: failedIDs,
	}
}

