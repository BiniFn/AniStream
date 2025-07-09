package models

import "github.com/coeeter/aniways/internal/repository"

type AnimeDto struct {
	ID          string `json:"id"`
	Ename       string `json:"ename"`
	JName       string `json:"jname"`
	ImageURL    string `json:"imageUrl"`
	Genre       string `json:"genre"`
	MalID       int32  `json:"malId"`
	AnilistID   int32  `json:"anilistId"`
	LastEpisode int32  `json:"lastEpisode"`
}

func (a AnimeDto) FromRepository(anime repository.Anime) AnimeDto {
	return AnimeDto{
		ID:          anime.ID,
		Ename:       anime.Ename,
		JName:       anime.Jname,
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		MalID:       anime.MalID.Int32,
		AnilistID:   anime.AnilistID.Int32,
		LastEpisode: anime.LastEpisode,
	}
}

type AnimeMetadataDto struct {
	MalID              int32   `json:"malId"`
	Description        string  `json:"description"`
	MainPictureURL     string  `json:"mainPictureUrl"`
	MediaType          string  `json:"mediaType"`
	Rating             string  `json:"rating"`
	AiringStatus       string  `json:"airingStatus"`
	AvgEpisodeDuration int32   `json:"avgEpisodeDuration"`
	TotalEpisodes      int32   `json:"totalEpisodes"`
	Studio             string  `json:"studio"`
	Rank               int32   `json:"rank"`
	Mean               float64 `json:"mean"`
	ScoringUsers       int32   `json:"scoringUsers"`
	Popularity         int32   `json:"popularity"`
	AiringStartDate    string  `json:"airingStartDate"`
	AiringEndDate      string  `json:"airingEndDate"`
	Source             string  `json:"source"`
	SeasonYear         int32   `json:"seasonYear"`
	Season             string  `json:"season"`
	TrailerEmbedURL    string  `json:"trailerEmbedUrl"`
}

func (m AnimeMetadataDto) FromRepository(metadata repository.AnimeMetadatum) AnimeMetadataDto {
	return AnimeMetadataDto{
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

type AnimeWithMetadataDto struct {
	ID          string           `json:"id"`
	Ename       string           `json:"ename"`
	JName       string           `json:"jname"`
	ImageURL    string           `json:"imageUrl"`
	Genre       string           `json:"genre"`
	MalID       int32            `json:"malId"`
	AnilistID   int32            `json:"anilistId"`
	LastEpisode int32            `json:"lastEpisode"`
	Metadata    AnimeMetadataDto `json:"metadata"`
}

func (a AnimeWithMetadataDto) FromRepository(anime repository.Anime, metadata repository.AnimeMetadatum) AnimeWithMetadataDto {
	return AnimeWithMetadataDto{
		ID:          anime.ID,
		Ename:       anime.Ename,
		JName:       anime.Jname,
		ImageURL:    anime.ImageUrl,
		Genre:       anime.Genre,
		MalID:       anime.MalID.Int32,
		AnilistID:   anime.AnilistID.Int32,
		LastEpisode: anime.LastEpisode,
		Metadata:    AnimeMetadataDto{}.FromRepository(metadata),
	}
}
