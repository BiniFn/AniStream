package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/coeeter/aniways/internal/infra/client/hianime"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	retryCount = 3
	retryDelay = 500 * time.Millisecond
)

type AdminService struct {
	repo       *repository.Queries
	scraper    *hianime.HianimeScraper
	jobManager *JobManager
}

func NewAdminService(repo *repository.Queries, scraper *hianime.HianimeScraper) *AdminService {
	return &AdminService{
		repo:       repo,
		scraper:    scraper,
		jobManager: NewJobManager(),
	}
}

func (s *AdminService) processAnimeDetails(ctx context.Context, hiAnimeID string) models.ReprocessResult {
	var lastErr error
	for attempt := 1; attempt <= retryCount; attempt++ {
		info, err := s.scraper.GetAnimeInfoByHiAnimeID(ctx, hiAnimeID)
		if err == nil {
			return s.updateOrCreateAnimeInDatabase(ctx, hiAnimeID, info)
		}
		lastErr = err

		if attempt < retryCount {
			time.Sleep(retryDelay)
		}
	}

	return models.ReprocessResult{
		HiAnimeID: hiAnimeID,
		Success:   false,
		Message:   fmt.Sprintf("Failed to fetch anime details after %d attempts: %v", retryCount, lastErr),
	}
}

func (s *AdminService) updateOrCreateAnimeInDatabase(ctx context.Context, hiAnimeID string, info hianime.ScrapedAnimeInfoDto) models.ReprocessResult {
	existingAnime, err := s.repo.GetAnimeByHiAnimeId(ctx, hiAnimeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return s.createNewAnime(ctx, hiAnimeID, info)
		}

		return models.ReprocessResult{
			HiAnimeID: hiAnimeID,
			Success:   false,
			Message:   fmt.Sprintf("Failed to check existing anime: %v", err),
		}
	}

	if s.hasChanges(existingAnime, info) {
		updateParams := repository.UpdateAnimeParams{
			ID:          existingAnime.ID,
			Ename:       info.EName,
			Jname:       info.JName,
			ImageUrl:    info.PosterURL,
			Genre:       info.Genre,
			HiAnimeID:   info.HiAnimeID,
			MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
			AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
			LastEpisode: int32(info.LastEpisode),
			Season:      repository.Season(strings.ToLower(info.Season)),
			SeasonYear:  int32(info.SeasonYear),
			UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		}

		if err := s.repo.UpdateAnime(ctx, updateParams); err != nil {
			return models.ReprocessResult{
				HiAnimeID: hiAnimeID,
				Success:   false,
				Message:   fmt.Sprintf("Failed to update anime in database: %v", err),
			}
		}

		return models.ReprocessResult{
			HiAnimeID: hiAnimeID,
			Success:   true,
			Message:   "Successfully updated anime details",
		}
	} else {
		return models.ReprocessResult{
			HiAnimeID: hiAnimeID,
			Success:   true,
			Message:   "No changes detected, skipped update",
		}
	}
}

func (s *AdminService) createNewAnime(ctx context.Context, hiAnimeID string, info hianime.ScrapedAnimeInfoDto) models.ReprocessResult {
	insertParams := repository.InsertAnimeParams{
		Ename:       info.EName,
		Jname:       info.JName,
		ImageUrl:    info.PosterURL,
		Genre:       info.Genre,
		HiAnimeID:   info.HiAnimeID,
		MalID:       pgtype.Int4{Int32: int32(info.MalID), Valid: info.MalID > 0},
		AnilistID:   pgtype.Int4{Int32: int32(info.AnilistID), Valid: info.AnilistID > 0},
		LastEpisode: int32(info.LastEpisode),
		Season:      repository.Season(strings.ToLower(info.Season)),
		SeasonYear:  int32(info.SeasonYear),
		CreatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	if err := s.repo.InsertAnime(ctx, insertParams); err != nil {
		return models.ReprocessResult{
			HiAnimeID: hiAnimeID,
			Success:   false,
			Message:   fmt.Sprintf("Failed to create new anime in database: %v", err),
		}
	}

	return models.ReprocessResult{
		HiAnimeID: hiAnimeID,
		Success:   true,
		Message:   "Successfully created new anime",
	}
}

func (s *AdminService) hasChanges(existing repository.Anime, newInfo hianime.ScrapedAnimeInfoDto) bool {
	if existing.Ename != newInfo.EName {
		return true
	}
	if existing.Jname != newInfo.JName {
		return true
	}
	if existing.ImageUrl != newInfo.PosterURL {
		return true
	}
	if existing.Genre != newInfo.Genre {
		return true
	}
	if existing.HiAnimeID != newInfo.HiAnimeID {
		return true
	}
	if existing.LastEpisode != int32(newInfo.LastEpisode) {
		return true
	}
	if existing.Season != repository.Season(strings.ToLower(newInfo.Season)) {
		return true
	}
	if existing.SeasonYear != int32(newInfo.SeasonYear) {
		return true
	}

	existingMalID := int32(0)
	if existing.MalID.Valid {
		existingMalID = existing.MalID.Int32
	}
	if existingMalID != int32(newInfo.MalID) {
		return true
	}

	existingAnilistID := int32(0)
	if existing.AnilistID.Valid {
		existingAnilistID = existing.AnilistID.Int32
	}
	if existingAnilistID != int32(newInfo.AnilistID) {
		return true
	}

	return false
}

func (s *AdminService) StartBulkReprocessFromFile(ctx context.Context, fileReader io.Reader) (models.BulkReprocessResponse, error) {
	job := s.jobManager.CreateJob([]string{})

	backgroundCtx := context.Background()
	go s.processBulkJobFromFile(backgroundCtx, job.ID, fileReader)

	return models.BulkReprocessResponse{
		JobID:   job.ID,
		Message: "Bulk reprocessing job started from file",
	}, nil
}

func (s *AdminService) processBulkJobFromFile(ctx context.Context, jobID string, fileReader io.Reader) {
	fileContent, err := io.ReadAll(fileReader)
	if err != nil {
		s.jobManager.FailJob(jobID, fmt.Errorf("error reading file: %w", err))
		return
	}

	total, err := s.countAndValidateAnimeIDs(strings.NewReader(string(fileContent)))
	if err != nil {
		s.jobManager.FailJob(jobID, err)
		return
	}

	s.jobManager.StartJob(jobID)
	s.jobManager.UpdateJobTotal(jobID, total)

	s.processAnimeIDsInChunks(ctx, jobID, strings.NewReader(string(fileContent)), total)
}

func (s *AdminService) countAndValidateAnimeIDs(fileReader io.Reader) (int, error) {
	decoder := json.NewDecoder(fileReader)

	token, err := decoder.Token()
	if err != nil {
		return 0, fmt.Errorf("invalid JSON format: %w", err)
	}

	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return 0, fmt.Errorf("expected JSON array")
	}

	total := 0

	for decoder.More() {
		var hiAnimeID string
		if err := decoder.Decode(&hiAnimeID); err != nil {
			return 0, fmt.Errorf("invalid anime ID at index %d: %w", total, err)
		}

		if hiAnimeID == "" {
			return 0, fmt.Errorf("empty anime ID at index %d", total)
		}

		total++

		if total > 10000 {
			return 0, fmt.Errorf("too many IDs requested, maximum allowed: 10000")
		}
	}

	token, err = decoder.Token()
	if err != nil || token != json.Delim(']') {
		return 0, fmt.Errorf("invalid JSON array format")
	}

	if total == 0 {
		return 0, fmt.Errorf("no anime IDs found in JSON array")
	}

	return total, nil
}

func (s *AdminService) processAnimeIDsInChunks(ctx context.Context, jobID string, fileReader io.Reader, total int) {
	decoder := json.NewDecoder(fileReader)

	token, err := decoder.Token()
	if err != nil {
		s.jobManager.FailJob(jobID, fmt.Errorf("invalid JSON format: %w", err))
		return
	}

	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		s.jobManager.FailJob(jobID, fmt.Errorf("expected JSON array"))
		return
	}

	chunkSize := 100
	maxWorkers := 10
	processed := 0

	var successResults []models.ReprocessResult
	var errorResults []models.ReprocessError
	var currentChunk []string

	for decoder.More() {
		var hiAnimeID string
		if err := decoder.Decode(&hiAnimeID); err != nil {
			s.jobManager.FailJob(jobID, fmt.Errorf("invalid anime ID: %w", err))
			return
		}

		currentChunk = append(currentChunk, hiAnimeID)

		if len(currentChunk) >= chunkSize {
			chunkResults, chunkErrors := s.processChunk(ctx, currentChunk, maxWorkers)

			processed += len(currentChunk)
			successResults = append(successResults, chunkResults...)
			errorResults = append(errorResults, chunkErrors...)

			s.jobManager.UpdateJobProgress(jobID, processed, len(successResults), len(errorResults),
				fmt.Sprintf("Processed %d/%d anime IDs", processed, total))

			currentChunk = nil

			select {
			case <-ctx.Done():
				s.jobManager.FailJob(jobID, ctx.Err())
				return
			default:
			}
		}
	}

	if len(currentChunk) > 0 {
		chunkResults, chunkErrors := s.processChunk(ctx, currentChunk, maxWorkers)

		processed += len(currentChunk)
		successResults = append(successResults, chunkResults...)
		errorResults = append(errorResults, chunkErrors...)
	}

	token, err = decoder.Token()
	if err != nil || token != json.Delim(']') {
		s.jobManager.FailJob(jobID, fmt.Errorf("invalid JSON array format"))
		return
	}

	s.jobManager.CompleteJob(jobID, successResults, errorResults)
}

func (s *AdminService) processChunk(ctx context.Context, chunk []string, maxWorkers int) ([]models.ReprocessResult, []models.ReprocessError) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var successResults []models.ReprocessResult
	var errorResults []models.ReprocessError

	semaphore := make(chan struct{}, maxWorkers)

	for _, hiAnimeID := range chunk {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := s.processAnimeDetails(ctx, id)

			mu.Lock()
			if result.Success {
				successResults = append(successResults, result)
			} else {
				errorResults = append(errorResults, models.ReprocessError{
					HiAnimeID: result.HiAnimeID,
					Error:     result.Message,
				})
			}
			mu.Unlock()
		}(hiAnimeID)
	}

	wg.Wait()
	return successResults, errorResults
}

func (s *AdminService) GetBulkJobStatus(jobID string) (models.BulkJobStatus, error) {
	job, exists := s.jobManager.GetJob(jobID)
	if !exists {
		return models.BulkJobStatus{}, fmt.Errorf("job not found")
	}

	return job.ToStatusResponse(), nil
}

func (s *AdminService) GetBulkJobResult(jobID string) (models.BulkJobResult, error) {
	job, exists := s.jobManager.GetJob(jobID)
	if !exists {
		return models.BulkJobResult{}, fmt.Errorf("job not found")
	}

	return job.ToResultResponse(), nil
}

func (s *AdminService) StartBulkReprocessFromIDs(hiAnimeIDs []string) (string, error) {
	if len(hiAnimeIDs) == 0 {
		return "", fmt.Errorf("no anime IDs provided")
	}

	if len(hiAnimeIDs) > 10000 {
		return "", fmt.Errorf("too many anime IDs, maximum 10,000 allowed")
	}

	job := s.jobManager.CreateJob(hiAnimeIDs)
	go s.processBulkJobFromIDs(job.ID, hiAnimeIDs)

	return job.ID, nil
}

func (s *AdminService) StartBulkReprocessFromIDsWithParent(hiAnimeIDs []string, parentJobID string) (string, error) {
	if len(hiAnimeIDs) == 0 {
		return "", fmt.Errorf("no anime IDs provided")
	}

	if len(hiAnimeIDs) > 10000 {
		return "", fmt.Errorf("too many anime IDs, maximum 10,000 allowed")
	}

	parentJob, exists := s.jobManager.GetJob(parentJobID)
	if !exists {
		return "", fmt.Errorf("parent job not found")
	}

	parentJob.Mu.RLock()
	parentSuccessResults := make([]models.ReprocessResult, len(parentJob.Results))
	copy(parentSuccessResults, parentJob.Results)
	parentJob.Mu.RUnlock()

	job := s.jobManager.CreateJob(hiAnimeIDs)
	job.Mu.Lock()
	job.Results = parentSuccessResults
	job.Success = len(parentSuccessResults)
	job.Mu.Unlock()

	go s.processBulkJobFromIDsWithParent(job.ID, hiAnimeIDs, parentSuccessResults)

	return job.ID, nil
}

func (s *AdminService) processBulkJobFromIDs(jobID string, hiAnimeIDs []string) {
	ctx := context.Background()
	job, exists := s.jobManager.GetJob(jobID)
	if !exists {
		return
	}

	var successResults []models.ReprocessResult
	var errorResults []models.ReprocessError

	chunks := s.chunkSlice(hiAnimeIDs, 100)

	for _, chunk := range chunks {
		success, errors := s.processChunk(ctx, chunk, 10)
		successResults = append(successResults, success...)
		errorResults = append(errorResults, errors...)

		job.Mu.Lock()
		job.Processed += len(chunk)
		job.Success = len(successResults)
		job.Failed = len(errorResults)
		job.Results = successResults
		job.Errors = errorResults
		job.Mu.Unlock()
	}

	job.Mu.Lock()
	if len(errorResults) == 0 {
		job.Status = models.JobStatusCompleted
		job.Message = "All anime processed successfully"
	} else {
		job.Status = models.JobStatusCompleted
		job.Message = fmt.Sprintf("Job completed with %d failures", len(errorResults))
	}
	job.Mu.Unlock()
}

func (s *AdminService) processBulkJobFromIDsWithParent(jobID string, hiAnimeIDs []string, parentResults []models.ReprocessResult) {
	ctx := context.Background()
	job, exists := s.jobManager.GetJob(jobID)
	if !exists {
		return
	}

	successResults := make([]models.ReprocessResult, len(parentResults))
	copy(successResults, parentResults)
	var errorResults []models.ReprocessError

	chunks := s.chunkSlice(hiAnimeIDs, 100)

	for _, chunk := range chunks {
		success, errors := s.processChunk(ctx, chunk, 10)
		successResults = append(successResults, success...)
		errorResults = append(errorResults, errors...)

		job.Mu.Lock()
		job.Processed += len(chunk)
		job.Success = len(successResults)
		job.Failed = len(errorResults)
		job.Results = successResults
		job.Errors = errorResults
		job.Mu.Unlock()
	}

	job.Mu.Lock()
	if len(errorResults) == 0 {
		job.Status = models.JobStatusCompleted
		job.Message = "All anime processed successfully"
	} else {
		job.Status = models.JobStatusCompleted
		job.Message = fmt.Sprintf("Job completed with %d failures", len(errorResults))
	}
	job.Mu.Unlock()
}

func (s *AdminService) chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := min(i+chunkSize, len(slice))
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
