package models

type AnimeResponse struct {
	ID          string  `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	EName       *string `json:"ename" example:"Attack on Titan"`
	JName       *string `json:"jname" example:"進撃の巨人"`
	ImageURL    string  `json:"imageUrl" validate:"required" example:"https://example.com/anime/image.jpg"`
	Genre       string  `json:"genre" validate:"required" example:"Action, Drama"`
	Season      string  `json:"season" validate:"required" example:"spring"`
	SeasonYear  *int32  `json:"seasonYear" example:"2023"`
	MalID       *int32  `json:"malId" example:"12345"`
	AnilistID   *int32  `json:"anilistId" example:"67890"`
	LastEpisode *int32  `json:"lastEpisode" example:"25"`
}

type AnimeWithMetadataResponse struct {
	ID          string                 `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	EName       *string                `json:"ename" example:"Attack on Titan"`
	JName       *string                `json:"jname" example:"進撃の巨人"`
	ImageURL    string                 `json:"imageUrl" validate:"required" example:"https://example.com/anime/image.jpg"`
	Genre       string                 `json:"genre" validate:"required" example:"Action, Drama"`
	Season      string                 `json:"season" validate:"required" example:"spring"`
	SeasonYear  *int32                 `json:"seasonYear" example:"2023"`
	MalID       *int32                 `json:"malId" example:"12345"`
	AnilistID   *int32                 `json:"anilistId" example:"67890"`
	LastEpisode *int32                 `json:"lastEpisode" example:"25"`
	Metadata    *AnimeMetadataResponse `json:"metadata"`
}

type AnimeMetadataResponse struct {
	MalID              int32   `json:"malId" validate:"required" example:"12345"`
	Description        string  `json:"description" validate:"required" example:"Humanity fights for survival against giant humanoid Titans."`
	MainPictureURL     string  `json:"mainPictureUrl" validate:"required" example:"https://example.com/main.jpg"`
	MediaType          string  `json:"mediaType" validate:"required" example:"TV"`
	Rating             string  `json:"rating" validate:"required" example:"R - 17+ (violence & profanity)"`
	AiringStatus       string  `json:"airingStatus" validate:"required" example:"finished_airing"`
	AvgEpisodeDuration int32   `json:"avgEpisodeDuration" validate:"required" example:"24"`
	TotalEpisodes      int32   `json:"totalEpisodes" validate:"required" example:"25"`
	Studio             string  `json:"studio" validate:"required" example:"Studio Pierrot"`
	Rank               int32   `json:"rank" validate:"required" example:"5"`
	Mean               float64 `json:"mean" validate:"required" example:"9.0"`
	ScoringUsers       int32   `json:"scoringUsers" validate:"required" example:"500000"`
	Popularity         int32   `json:"popularity" validate:"required" example:"1"`
	AiringStartDate    string  `json:"airingStartDate" validate:"required" example:"2023-04-01"`
	AiringEndDate      string  `json:"airingEndDate" validate:"required" example:"2023-09-30"`
	Source             string  `json:"source" validate:"required" example:"Manga"`
	SeasonYear         int32   `json:"seasonYear" validate:"required" example:"2023"`
	Season             string  `json:"season" validate:"required" example:"spring"`
	TrailerEmbedURL    string  `json:"trailerEmbedUrl" validate:"required" example:"https://www.youtube.com/embed/abc123"`
}

type TrailerResponse struct {
	Trailer string `json:"trailer" validate:"required" example:"https://www.youtube.com/embed/abc123"`
}

type BannerResponse struct {
	URL string `json:"url" validate:"required" example:"https://example.com/banner.jpg"`
}

type RelationsResponse struct {
	WatchOrder []AnimeResponse `json:"watchOrder" validate:"required"`
	Related    []AnimeResponse `json:"related" validate:"required"`
}

type SeasonalAnimeResponse struct {
	ID             string        `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	BannerImageURL string        `json:"bannerImageUrl" validate:"required" example:"https://example.com/banner.jpg"`
	Description    string        `json:"description" validate:"required" example:"New seasonal anime description"`
	StartDate      int64         `json:"startDate" validate:"required" example:"1672531200"`
	Type           string        `json:"type" validate:"required" example:"TV"`
	Episodes       int32         `json:"episodes" validate:"required" example:"12"`
	Anime          AnimeResponse `json:"anime" validate:"required"`
}

type EpisodeResponse struct {
	ID       string `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	Title    string `json:"title" validate:"required" example:"The Attack Titan"`
	Number   int    `json:"number" validate:"required" example:"1"`
	IsFiller bool   `json:"isFiller" validate:"required" example:"false"`
}

type EpisodeServerResponse struct {
	Type       string `json:"type" validate:"required" example:"sub"`
	ServerName string `json:"serverName" validate:"required" example:"vidstreaming"`
	ServerID   string `json:"serverId" validate:"required" example:"V1StGXR8Z5jdHi6B"`
}

type StreamingDataResponse struct {
	Source StreamingSourceResponse `json:"source" validate:"required"`
	Intro  SegmentResponse         `json:"intro" validate:"required"`
	Outro  SegmentResponse         `json:"outro" validate:"required"`
	Tracks []TrackResponse         `json:"tracks" validate:"required"`
}

type StreamingSourceResponse struct {
	Hls      *string `json:"hls" example:"https://example.com/stream.m3u8"`
	ProxyHls *string `json:"proxyHls" example:"/proxy?p=encodedurl&s=hd"`
	Iframe   string  `json:"iframe" validate:"required" example:"https://example.com/embed/abc123"`
}

type SegmentResponse struct {
	Start int `json:"start" validate:"required" example:"90"`
	End   int `json:"end" validate:"required" example:"180"`
}

type TrackResponse struct {
	URL     string `json:"url" validate:"required" example:"/proxy?p=encodedurl&s=hd"`
	Raw     string `json:"raw" validate:"required" example:"https://example.com/subtitles.vtt"`
	Kind    string `json:"kind" validate:"required" example:"captions"`
	Label   string `json:"label" validate:"required" example:"English"`
	Default bool   `json:"default" validate:"required" example:"true"`
}

type AnimeWithLibraryResponse struct {
	ID          string       `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	EName       *string      `json:"ename" example:"Attack on Titan"`
	JName       *string      `json:"jname" example:"進撃の巨人"`
	ImageURL    string       `json:"imageUrl" validate:"required" example:"https://example.com/anime/image.jpg"`
	Genre       string       `json:"genre" validate:"required" example:"Action, Drama"`
	Season      string       `json:"season" validate:"required" example:"spring"`
	SeasonYear  *int32       `json:"seasonYear" example:"2023"`
	MalID       *int32       `json:"malId" example:"12345"`
	AnilistID   *int32       `json:"anilistId" example:"67890"`
	LastEpisode *int32       `json:"lastEpisode" example:"25"`
	Library     *LibraryInfo `json:"library,omitempty"`
}

type LibraryInfo struct {
	ID              string        `json:"id" validate:"required" example:"V1StGXR8Z5jdHi6B"`
	Status          LibraryStatus `json:"status" validate:"required" example:"watching"`
	WatchedEpisodes int32         `json:"watchedEpisodes" validate:"required,min=0" example:"12"`
	UpdatedAt       string        `json:"updatedAt" validate:"required" example:"2023-01-01T00:00:00Z"`
}

type AnimeListResponse = Pagination[AnimeResponse]
type AnimeWithLibraryListResponse = Pagination[AnimeWithLibraryResponse]
type SeasonalAnimeListResponse = []SeasonalAnimeResponse
type EpisodeListResponse = []EpisodeResponse
type EpisodeServerListResponse = []EpisodeServerResponse
type TrendingAnimeListResponse = []AnimeResponse
type PopularAnimeListResponse = []AnimeResponse
