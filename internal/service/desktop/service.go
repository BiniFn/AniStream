package desktop

import (
	"context"

	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type DesktopService struct {
	repo *repository.Queries
}

func NewDesktopService(repo *repository.Queries) *DesktopService {
	return &DesktopService{
		repo: repo,
	}
}

func (s *DesktopService) GetAllReleases(ctx context.Context) ([]models.DesktopReleaseResponse, error) {
	releases, err := s.repo.GetAllDesktopReleases(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.DesktopReleaseResponse, 0, len(releases))
	for _, r := range releases {
		result = append(result, mapRelease(r))
	}

	return result, nil
}

func (s *DesktopService) GetLatestReleases(ctx context.Context) (*models.DesktopVersionResponse, error) {
	releases, err := s.repo.GetLatestDesktopReleases(ctx)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, nil
	}

	return groupReleasesByVersion(releases), nil
}

func (s *DesktopService) GetReleasesByVersion(ctx context.Context, version string) (*models.DesktopVersionResponse, error) {
	releases, err := s.repo.GetDesktopReleasesByVersion(ctx, version)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, nil
	}

	return groupReleasesByVersion(releases), nil
}

type CreateReleaseParams struct {
	Version      string
	Platform     string
	DownloadURL  string
	FileName     string
	FileSize     int64
	ReleaseNotes string
}

func (s *DesktopService) CreateRelease(ctx context.Context, params CreateReleaseParams) (models.DesktopReleaseResponse, error) {
	var releaseNotes pgtype.Text
	if params.ReleaseNotes != "" {
		releaseNotes = pgtype.Text{String: params.ReleaseNotes, Valid: true}
	}

	release, err := s.repo.InsertDesktopRelease(ctx, repository.InsertDesktopReleaseParams{
		Version:      params.Version,
		Platform:     repository.DesktopPlatform(params.Platform),
		DownloadUrl:  params.DownloadURL,
		FileName:     params.FileName,
		FileSize:     params.FileSize,
		ReleaseNotes: releaseNotes,
	})
	if err != nil {
		return models.DesktopReleaseResponse{}, err
	}

	return mapRelease(release), nil
}

func (s *DesktopService) DeleteReleasesByVersion(ctx context.Context, version string) error {
	return s.repo.DeleteDesktopReleasesByVersion(ctx, version)
}

func mapRelease(r repository.DesktopRelease) models.DesktopReleaseResponse {
	var releaseNotes *string
	if r.ReleaseNotes.Valid {
		releaseNotes = &r.ReleaseNotes.String
	}

	return models.DesktopReleaseResponse{
		ID:           r.ID,
		Version:      r.Version,
		Platform:     string(r.Platform),
		DownloadURL:  r.DownloadUrl,
		FileName:     r.FileName,
		FileSize:     r.FileSize,
		ReleaseNotes: releaseNotes,
		CreatedAt:    r.CreatedAt.Time,
	}
}

func groupReleasesByVersion(releases []repository.DesktopRelease) *models.DesktopVersionResponse {
	if len(releases) == 0 {
		return nil
	}

	var releaseNotes *string
	if releases[0].ReleaseNotes.Valid {
		releaseNotes = &releases[0].ReleaseNotes.String
	}

	platforms := make([]models.DesktopPlatformRelease, 0, len(releases))
	for _, r := range releases {
		platforms = append(platforms, models.DesktopPlatformRelease{
			Platform:    string(r.Platform),
			DownloadURL: r.DownloadUrl,
			FileName:    r.FileName,
			FileSize:    r.FileSize,
		})
	}

	return &models.DesktopVersionResponse{
		Version:      releases[0].Version,
		ReleaseNotes: releaseNotes,
		CreatedAt:    releases[0].CreatedAt.Time,
		Platforms:    platforms,
	}
}
