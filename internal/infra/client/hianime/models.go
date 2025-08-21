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
	MalID       int    `json:"malId"`
	AnilistID   int    `json:"anilistId"`
	LastEpisode int    `json:"lastEpisode"`
	Season      string `json:"season"`
	SeasonYear  int    `json:"seasonYear"`
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

type ScrapedSegment struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type ScrapedTrack struct {
	File    string `json:"file"`
	Kind    string `json:"kind"`
	Label   string `json:"label,omitempty"`
	Default bool   `json:"default,omitempty"`
}

type ScrapedEpisodeSourceDto struct {
	Hls    *string `json:"hls"`
	Iframe string  `json:"iframe"`
}

type ScrapedStreamData struct {
	Source ScrapedEpisodeSourceDto `json:"source"`
	Intro  ScrapedSegment          `json:"intro"`
	Outro  ScrapedSegment          `json:"outro"`
	Tracks []ScrapedTrack          `json:"tracks"`
	Server string                  `json:"server"`
}
