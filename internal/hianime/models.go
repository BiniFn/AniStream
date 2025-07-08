package hianime

type PageInfo struct {
	TotalPages      int  `json:"totalPages"`
	CurrentPage     int  `json:"currentPage"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

type Pagination[T any] struct {
	PageInfo PageInfo `json:"pageInfo"`
	Items    []T      `json:"items"`
}

type ScrapedAnimeInfoDto struct {
	HiAnimeID   string `json:"hiAnimeId"`
	EName       string `json:"eName"`
	JName       string `json:"jName"`
	PosterURL   string `json:"posterUrl"`
	Genre       string `json:"genre"`
	MalID       string `json:"malId"`
	AnilistID   string `json:"anilistId"`
	LastEpisode int    `json:"lastEpisode"`
}

type ScrapedEpisodeDto struct {
	EpisodeID string `json:"episodeId"`
	Title     string `json:"title"`
	Number    int    `json:"number"`
	IsFiller  bool   `json:"isFiller"`
}

type ScrapedEpisodeServerDto struct {
	Type       string `json:"type"`
	ServerName string `json:"serverName"`
	ServerID   string `json:"serverId"`
}
