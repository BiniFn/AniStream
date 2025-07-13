package models

import (
	"sort"

	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/client/hianime/streaming"
	"github.com/coeeter/aniways/internal/repository"
)

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

func (a AnimeDto) FromSearch(anime repository.SearchAnimesRow) AnimeDto {
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

type TrailerDto struct {
	Trailer string `json:"trailer"`
}

type EpisodeDto struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Number   int    `json:"number"`
	IsFiller bool   `json:"isFiller"`
}

func (e EpisodeDto) FromScraper(episode hianime.ScrapedEpisodeDto) EpisodeDto {
	return EpisodeDto{
		ID:       episode.EpisodeID,
		Title:    episode.Title,
		Number:   episode.Number,
		IsFiller: episode.IsFiller,
	}
}

type IndividualServerDto struct {
	Type       string `json:"type"`
	ServerName string `json:"serverName"`
	ServerID   string `json:"serverId"`
}

type ServerDto struct {
	Sub []IndividualServerDto `json:"sub"`
	Dub []IndividualServerDto `json:"dub"`
	Raw []IndividualServerDto `json:"raw"`
}

func sortByName(s []IndividualServerDto) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].ServerName < s[j].ServerName
	})
}

func (s ServerDto) FromScraper(servers []hianime.ScrapedEpisodeServerDto) ServerDto {
	sub, dub, raw := []IndividualServerDto{}, []IndividualServerDto{}, []IndividualServerDto{}

	for _, server := range servers {
		dto := IndividualServerDto{
			Type:       server.Type,
			ServerName: server.ServerName,
			ServerID:   server.ServerID,
		}
		switch server.Type {
		case "sub":
			sub = append(sub, dto)
		case "dub":
			dub = append(dub, dto)
		case "raw":
			raw = append(raw, dto)
		}
	}

	sortByName(sub)
	sortByName(dub)
	sortByName(raw)

	return ServerDto{
		Sub: sub,
		Dub: dub,
		Raw: raw,
	}
}

type StreamingDataDto struct {
	Source     string     `json:"source"`
	ServerName string     `json:"serverName"`
	Type       string     `json:"type"`
	Intro      SegmentDto `json:"intro"`
	Outro      SegmentDto `json:"outro"`
	Tracks     []TrackDto `json:"tracks"`
}

type SegmentDto struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type TrackDto struct {
	File    string `json:"file"`
	Kind    string `json:"kind"`
	Label   string `json:"label,omitempty"`
	Default bool   `json:"default,omitempty"`
}

func (s StreamingDataDto) FromScraper(data streaming.ScrapedUnencryptedSources) StreamingDataDto {
	return StreamingDataDto{
		Source:     data.Source,
		ServerName: data.ServerName,
		Type:       data.Type,
		Intro: SegmentDto{
			Start: data.Intro.Start,
			End:   data.Intro.End,
		},
		Outro: SegmentDto{
			Start: data.Outro.Start,
			End:   data.Outro.End,
		},
		Tracks: func(t []streaming.ScrapedTrack) []TrackDto {
			tracks := make([]TrackDto, len(t))
			for i, track := range t {
				tracks[i] = TrackDto{
					File:    track.File,
					Kind:    track.Kind,
					Label:   track.Label,
					Default: track.Default,
				}
			}
			return tracks
		}(data.Tracks),
	}
}

type SeasonalAnimeDto struct {
	ID             string   `json:"id"`
	BannerImageURL string   `json:"bannerImageUrl"`
	Description    string   `json:"description"`
	StartDate      int64    `json:"startDate"`
	Type           string   `json:"type"`
	Episodes       int32    `json:"episodes"`
	Anime          AnimeDto `json:"anime"`
}
