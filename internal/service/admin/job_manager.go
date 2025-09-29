package admin

import (
	"fmt"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/models"
	"github.com/google/uuid"
)

type JobManager struct {
	jobs map[string]*models.BulkJob
	mu   sync.RWMutex
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]*models.BulkJob),
	}
}

func (jm *JobManager) CreateJob(hiAnimeIDs []string) *models.BulkJob {
	jobID := uuid.New().String()
	job := &models.BulkJob{
		ID:         jobID,
		Status:     models.JobStatusPending,
		Total:      len(hiAnimeIDs),
		Processed:  0,
		Success:    0,
		Failed:     0,
		Message:    "Job created, waiting to start",
		StartedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		HiAnimeIDs: hiAnimeIDs,
		Results:    make([]models.ReprocessResult, 0),
		Errors:     make([]models.ReprocessError, 0),
	}

	jm.mu.Lock()
	jm.jobs[jobID] = job
	jm.mu.Unlock()

	return job
}

func (jm *JobManager) GetJob(jobID string) (*models.BulkJob, bool) {
	jm.mu.RLock()
	job, exists := jm.jobs[jobID]
	jm.mu.RUnlock()
	return job, exists
}

func (jm *JobManager) UpdateJobProgress(jobID string, processed, success, failed int, message string) {
	jm.mu.RLock()
	job, exists := jm.jobs[jobID]
	jm.mu.RUnlock()

	if !exists {
		return
	}

	job.Mu.Lock()
	job.Processed = processed
	job.Success = success
	job.Failed = failed
	job.Message = message
	job.UpdatedAt = time.Now()
	job.Mu.Unlock()
}

func (jm *JobManager) CompleteJob(jobID string, results []models.ReprocessResult, errors []models.ReprocessError) {
	jm.mu.RLock()
	job, exists := jm.jobs[jobID]
	jm.mu.RUnlock()

	if !exists {
		return
	}

	job.Mu.Lock()
	job.Status = models.JobStatusCompleted
	job.Processed = job.Total
	job.Success = len(results)
	job.Failed = len(errors)
	job.Message = "Job completed successfully"
	job.UpdatedAt = time.Now()
	job.Results = results
	job.Errors = errors
	job.Mu.Unlock()
}

func (jm *JobManager) FailJob(jobID string, err error) {
	jm.mu.RLock()
	job, exists := jm.jobs[jobID]
	jm.mu.RUnlock()

	if !exists {
		return
	}

	job.Mu.Lock()
	job.Status = models.JobStatusFailed
	job.Message = fmt.Sprintf("Job failed: %v", err)
	job.UpdatedAt = time.Now()
	job.Mu.Unlock()
}

func (jm *JobManager) StartJob(jobID string) {
	jm.mu.RLock()
	job, exists := jm.jobs[jobID]
	jm.mu.RUnlock()

	if !exists {
		return
	}

	job.Mu.Lock()
	job.Status = models.JobStatusRunning
	job.Message = "Job started, processing anime details..."
	job.UpdatedAt = time.Now()
	job.Mu.Unlock()
}

func (jm *JobManager) UpdateJobTotal(jobID string, total int) {
	jm.mu.RLock()
	job, exists := jm.jobs[jobID]
	jm.mu.RUnlock()

	if !exists {
		return
	}

	job.Mu.Lock()
	job.Total = total
	job.Message = fmt.Sprintf("Job initialized with %d anime IDs", total)
	job.UpdatedAt = time.Now()
	job.Mu.Unlock()
}
