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
