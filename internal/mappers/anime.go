package mappers

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/coeeter/aniways/internal/infra/client/hianime"
	"github.com/coeeter/aniways/internal/models"
	"github.com/coeeter/aniways/internal/repository"
)

func AnimeFromRepository(anime repository.Anime) models.AnimeResponse {
	return models.AnimeResponse{
		ID:          anime.ID,
		EName:       nilIfEmpty(anime.Ename),
		JName:       nilIfEmpty(anime.Jname),
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		Season:      string(anime.Season),
		SeasonYear:  nilIfEmpty(anime.SeasonYear),
		MalID:       nilIfEmpty(anime.MalID.Int32),
		AnilistID:   nilIfEmpty(anime.AnilistID.Int32),
		LastEpisode: nilIfEmpty(anime.LastEpisode),
	}
}

func AnimeFromCatalog(anime repository.GetAnimeCatalogRow) models.AnimeResponse {
	return models.AnimeResponse{
		ID:          anime.ID,
		EName:       nilIfEmpty(anime.Ename),
		JName:       nilIfEmpty(anime.Jname),
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		Season:      string(anime.Season),
		SeasonYear:  nilIfEmpty(anime.SeasonYear),
		MalID:       nilIfEmpty(anime.MalID.Int32),
		AnilistID:   nilIfEmpty(anime.AnilistID.Int32),
		LastEpisode: nilIfEmpty(anime.LastEpisode),
	}
}

func AnimeWithLibraryFromCatalog(anime repository.GetAnimeCatalogRow) models.AnimeWithLibraryResponse {
	response := models.AnimeWithLibraryResponse{
		ID:          anime.ID,
		EName:       nilIfEmpty(anime.Ename),
		JName:       nilIfEmpty(anime.Jname),
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		Season:      string(anime.Season),
		SeasonYear:  nilIfEmpty(anime.SeasonYear),
		MalID:       nilIfEmpty(anime.MalID.Int32),
		AnilistID:   nilIfEmpty(anime.AnilistID.Int32),
		LastEpisode: nilIfEmpty(anime.LastEpisode),
	}

	if anime.LibraryID.Valid {
		library := models.LibraryInfo{
			ID:              anime.LibraryID.String,
			WatchedEpisodes: anime.LibraryWatchedEpisodes.Int32,
		}

		if anime.LibraryStatus.Valid {
			library.Status = models.LibraryStatus(anime.LibraryStatus.LibraryStatus)
		}

		if anime.LibraryUpdatedAt.Valid {
			library.UpdatedAt = anime.LibraryUpdatedAt.Time.Format(time.RFC3339)
		}

		response.Library = &library
	}

	return response
}

func AnimeFromSearch(anime repository.SearchAnimesRow) models.AnimeResponse {
	return models.AnimeResponse{
		ID:          anime.ID,
		EName:       nilIfEmpty(anime.Ename),
		JName:       nilIfEmpty(anime.Jname),
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		Season:      string(anime.Season),
		SeasonYear:  nilIfEmpty(anime.SeasonYear),
		MalID:       nilIfEmpty(anime.MalID.Int32),
		AnilistID:   nilIfEmpty(anime.AnilistID.Int32),
		LastEpisode: nilIfEmpty(anime.LastEpisode),
	}
}

func AnimeMetadataFromRepository(metadata repository.AnimeMetadatum) models.AnimeMetadataResponse {
	return models.AnimeMetadataResponse{
		MalID:              metadata.MalID,
		Description:        metadata.Description.String,
		MainPictureURL:     metadata.MainPictureUrl.String,
		MediaType:          metadata.MediaType.String,
		Rating:             string(metadata.Rating),
		AiringStatus:       string(metadata.AiringStatus),
		AvgEpisodeDuration: metadata.AvgEpisodeDuration.Int32,
		TotalEpisodes:      metadata.TotalEpisodes.Int32,
		Studio:             metadata.Studio.String,
		Rank:               metadata.Rank.Int32,
		Mean:               metadata.Mean.Float64,
		ScoringUsers:       metadata.Scoringusers.Int32,
		Popularity:         metadata.Popularity.Int32,
		AiringStartDate:    metadata.AiringStartDate.String,
		AiringEndDate:      metadata.AiringEndDate.String,
		Source:             metadata.Source.String,
		SeasonYear:         metadata.SeasonYear.Int32,
		Season:             string(metadata.Season),
		TrailerEmbedURL:    metadata.TrailerEmbedUrl.String,
	}
}

func AnimeWithMetadataFromRepository(anime repository.Anime, metadata repository.AnimeMetadatum) models.AnimeWithMetadataResponse {
	meta := AnimeMetadataFromRepository(metadata)

	var metaPointer *models.AnimeMetadataResponse
	if meta.MalID != 0 {
		metaPointer = &meta
	}

	return models.AnimeWithMetadataResponse{
		ID:          anime.ID,
		EName:       nilIfEmpty(anime.Ename),
		JName:       nilIfEmpty(anime.Jname),
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		Season:      string(anime.Season),
		SeasonYear:  nilIfEmpty(anime.SeasonYear),
		MalID:       nilIfEmpty(anime.MalID.Int32),
		AnilistID:   nilIfEmpty(anime.AnilistID.Int32),
		LastEpisode: nilIfEmpty(anime.LastEpisode),
		Metadata:    metaPointer,
	}
}

func EpisodeFromScraper(episode hianime.ScrapedEpisodeDto) models.EpisodeResponse {
	return models.EpisodeResponse{
		ID:       episode.EpisodeID,
		Title:    episode.Title,
		Number:   episode.Number,
		IsFiller: episode.IsFiller,
	}
}

func EpisodeServerFromScraper(server hianime.ScrapedEpisodeServerDto) models.EpisodeServerResponse {
	return models.EpisodeServerResponse{
		Type:       server.Type,
		ServerName: server.ServerName,
		ServerID:   server.ServerID,
	}
}

func StreamingDataFromScraper(data hianime.ScrapedStreamData) models.StreamingDataResponse {
	source := models.StreamingSourceResponse{
		Iframe: data.Source.Iframe,
	}
	if data.Source.Hls != nil {
		source.Hls = data.Source.Hls
		p := base64.StdEncoding.EncodeToString([]byte(*data.Source.Hls))
		proxy := fmt.Sprintf("/proxy?p=%s&s=hd", p)
		source.ProxyHls = &proxy
	}

	tracks := make([]models.TrackResponse, len(data.Tracks))
	for i, track := range data.Tracks {
		encoder := base64.StdEncoding
		p := encoder.EncodeToString([]byte(track.File))
		s := "hd"

		tracks[i] = models.TrackResponse{
			URL:     fmt.Sprintf("/proxy?p=%s&s=%s", p, s),
			Raw:     track.File,
			Kind:    track.Kind,
			Label:   track.Label,
			Default: track.Default,
		}
	}

	return models.StreamingDataResponse{
		Source: source,
		Intro: models.SegmentResponse{
			Start: data.Intro.Start,
			End:   data.Intro.End,
		},
		Outro: models.SegmentResponse{
			Start: data.Outro.Start,
			End:   data.Outro.End,
		},
		Tracks: tracks,
	}
}

func nilIfEmpty[T comparable](value T) *T {
	var zero T
	if value == zero {
		return nil
	}
	return &value
}
