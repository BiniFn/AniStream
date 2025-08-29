package mappers

import (
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

func LibraryFromRepository(l repository.Library, a repository.Anime) models.LibraryResponse {
	return models.LibraryResponse{
		ID:              l.ID,
		UserID:          l.UserID,
		AnimeID:         l.AnimeID,
		Status:          models.LibraryStatus(l.Status),
		WatchedEpisodes: l.WatchedEpisodes,
		CreatedAt:       l.CreatedAt.Time,
		UpdatedAt:       l.UpdatedAt.Time,
		Anime:           AnimeFromRepository(a),
	}
}

func LibraryImportJobFromRepository(j repository.LibraryImportJob) models.LibraryImportJobResponse {
	return models.LibraryImportJobResponse{
		ID:          j.ID,
		UserID:      j.UserID,
		Status:      models.LibraryStatus(j.Status),
		CreatedAt:   j.CreatedAt.Time,
		UpdatedAt:   j.UpdatedAt.Time,
		CompletedAt: j.CompletedAt.Time,
	}
}
