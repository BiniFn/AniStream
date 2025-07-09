package myanimelist

import (
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type MalAnimeMetadata struct {
	MalID       int    `json:"id"`
	Title       string `json:"title"`
	MainPicture struct {
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"main_picture"`
	AlternativeTitles struct {
		English  string   `json:"en"`
		Japanese string   `json:"ja"`
		Synonyms []string `json:"synonyms"`
	} `json:"alternative_titles"`
	Synopsis      string `json:"synopsis"`
	MediaType     string `json:"media_type"`
	Rating        string `json:"rating"`
	AvgEpDuration int    `json:"average_episode_duration"`
	Status        string `json:"status"`
	NumEpisodes   int    `json:"num_episodes"`
	Studios       []struct {
		MalID int    `json:"id"`
		Name  string `json:"name"`
	} `json:"studios"`
	Rank            int     `json:"rank"`
	Mean            float64 `json:"mean"`
	NumScoringUsers int     `json:"num_scoring_users"`
	Popularity      int     `json:"popularity"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	Source          string  `json:"source"`
	StartSeason     struct {
		Year   int    `json:"year"`
		Season string `json:"season"`
	} `json:"start_season"`
}

func (m MalAnimeMetadata) ToRepository() *repository.AnimeMetadatum {
	mainPictureUrl := m.MainPicture.Large
	if mainPictureUrl == "" {
		mainPictureUrl = m.MainPicture.Medium
	}

	rating := repository.RatingUnknown
	if m.Rating != "" {
		rating = repository.Rating(m.Rating)
	}

	airingStatus := repository.AiringStatusUnknown
	if m.Status != "" {
		airingStatus = repository.AiringStatus(m.Status)
	}

	season := repository.SeasonUnknown
	if m.StartSeason.Season != "" {
		season = repository.Season(m.StartSeason.Season)
	}

	return &repository.AnimeMetadatum{
		MalID:              int32(m.MalID),
		Description:        pgtype.Text{String: m.Synopsis, Valid: true},
		MainPictureUrl:     pgtype.Text{String: mainPictureUrl, Valid: true},
		MediaType:          pgtype.Text{String: m.MediaType, Valid: true},
		Rating:             rating,
		AiringStatus:       airingStatus,
		AvgEpisodeDuration: pgtype.Int4{Int32: int32(m.AvgEpDuration), Valid: true},
		TotalEpisodes:      pgtype.Int4{Int32: int32(m.NumEpisodes), Valid: true},
		Studio:             pgtype.Text{String: m.Studios[0].Name, Valid: true},
		Rank:               pgtype.Int4{Int32: int32(m.Rank), Valid: true},
		Mean:               pgtype.Float8{Float64: m.Mean, Valid: true},
		Scoringusers:       pgtype.Int4{Int32: int32(m.NumScoringUsers), Valid: true},
		Popularity:         pgtype.Int4{Int32: int32(m.Popularity), Valid: true},
		AiringStartDate:    pgtype.Text{String: m.StartDate, Valid: true},
		AiringEndDate:      pgtype.Text{String: m.EndDate, Valid: true},
		Source:             pgtype.Text{String: m.Source, Valid: true},
		SeasonYear:         pgtype.Int4{Int32: int32(m.StartSeason.Year), Valid: true},
		Season:             season,
	}
}
