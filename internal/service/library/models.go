package library

import (
	"time"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/coeeter/aniways/internal/service/anime"
)

type LibraryDto struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userId"`
	AnimeID         string         `json:"animeId"`
	Status          string         `json:"status"`
	WatchedEpisodes int32          `json:"watchedEpisodes"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	Anime           anime.AnimeDto `json:"anime"`
}

func (LibraryDto) FromRepository(l repository.Library, a repository.Anime) LibraryDto {
	return LibraryDto{
		ID:              l.ID,
		UserID:          l.UserID,
		AnimeID:         l.AnimeID,
		Status:          string(l.Status),
		WatchedEpisodes: l.WatchedEpisodes,
		CreatedAt:       l.CreatedAt.Time,
		UpdatedAt:       l.UpdatedAt.Time,
		Anime:           anime.AnimeDto{}.FromRepository(a),
	}
}
