package models

type AnimeResponse struct {
	ID          string  `json:"id" example:"V1StGXR8Z5jdHi6B"`
	EName       *string `json:"ename" example:"Attack on Titan"`
	JName       *string `json:"jname" example:"進撃の巨人"`
	ImageURL    string  `json:"imageUrl" example:"https://example.com/anime/image.jpg"`
	Genre       string  `json:"genre" example:"Action, Drama"`
	Season      string  `json:"season" example:"spring"`
	SeasonYear  *int32  `json:"seasonYear" example:"2023"`
	MalID       *int32  `json:"malId" example:"12345"`
	AnilistID   *int32  `json:"anilistId" example:"67890"`
	LastEpisode *int32  `json:"lastEpisode" example:"25"`
}

type AnimeWithMetadataResponse struct {
	ID          string                 `json:"id" example:"V1StGXR8Z5jdHi6B"`
	EName       *string                `json:"ename" example:"Attack on Titan"`
	JName       *string                `json:"jname" example:"進撃の巨人"`
	ImageURL    string                 `json:"imageUrl" example:"https://example.com/anime/image.jpg"`
	Genre       string                 `json:"genre" example:"Action, Drama"`
	Season      string                 `json:"season" example:"spring"`
	SeasonYear  *int32                 `json:"seasonYear" example:"2023"`
	MalID       *int32                 `json:"malId" example:"12345"`
	AnilistID   *int32                 `json:"anilistId" example:"67890"`
	LastEpisode *int32                 `json:"lastEpisode" example:"25"`
	Metadata    *AnimeMetadataResponse `json:"metadata"`
}

type AnimeMetadataResponse struct {
	MalID              int32   `json:"malId" example:"12345"`
	Description        string  `json:"description" example:"Humanity fights for survival against giant humanoid Titans."`
	MainPictureURL     string  `json:"mainPictureUrl" example:"https://example.com/main.jpg"`
	MediaType          string  `json:"mediaType" example:"TV"`
	Rating             string  `json:"rating" example:"R - 17+ (violence & profanity)"`
	AiringStatus       string  `json:"airingStatus" example:"finished_airing"`
	AvgEpisodeDuration int32   `json:"avgEpisodeDuration" example:"24"`
	TotalEpisodes      int32   `json:"totalEpisodes" example:"25"`
	Studio             string  `json:"studio" example:"Studio Pierrot"`
	Rank               int32   `json:"rank" example:"5"`
	Mean               float64 `json:"mean" example:"9.0"`
	ScoringUsers       int32   `json:"scoringUsers" example:"500000"`
	Popularity         int32   `json:"popularity" example:"1"`
	AiringStartDate    string  `json:"airingStartDate" example:"2023-04-01"`
	AiringEndDate      string  `json:"airingEndDate" example:"2023-09-30"`
	Source             string  `json:"source" example:"Manga"`
	SeasonYear         int32   `json:"seasonYear" example:"2023"`
	Season             string  `json:"season" example:"spring"`
	TrailerEmbedURL    string  `json:"trailerEmbedUrl" example:"https://www.youtube.com/embed/abc123"`
}

type TrailerResponse struct {
	Trailer string `json:"trailer" example:"https://www.youtube.com/embed/abc123"`
}

type BannerResponse struct {
	URL string `json:"url" example:"https://example.com/banner.jpg"`
}

type RelationsResponse struct {
	WatchOrder []AnimeResponse `json:"watchOrder"`
	Related    []AnimeResponse `json:"related"`
}

type SeasonalAnimeResponse struct {
	ID             string        `json:"id" example:"V1StGXR8Z5jdHi6B"`
	BannerImageURL string        `json:"bannerImageUrl" example:"https://example.com/banner.jpg"`
	Description    string        `json:"description" example:"New seasonal anime description"`
	StartDate      int64         `json:"startDate" example:"1672531200"`
	Type           string        `json:"type" example:"TV"`
	Episodes       int32         `json:"episodes" example:"12"`
	Anime          AnimeResponse `json:"anime"`
}

type EpisodeResponse struct {
	ID       string `json:"id" example:"V1StGXR8Z5jdHi6B"`
	Title    string `json:"title" example:"The Attack Titan"`
	Number   int    `json:"number" example:"1"`
	IsFiller bool   `json:"isFiller" example:"false"`
}

type EpisodeServerResponse struct {
	Type       string `json:"type" example:"sub"`
	ServerName string `json:"serverName" example:"vidstreaming"`
	ServerID   string `json:"serverId" example:"V1StGXR8Z5jdHi6B"`
}

type StreamingDataResponse struct {
	Source StreamingSourceResponse `json:"source"`
	Intro  SegmentResponse         `json:"intro"`
	Outro  SegmentResponse         `json:"outro"`
	Tracks []TrackResponse         `json:"tracks"`
}

type StreamingSourceResponse struct {
	Hls      *string `json:"hls" example:"https://example.com/stream.m3u8"`
	ProxyHls *string `json:"proxyHls" example:"/proxy?p=encodedurl&s=hd"`
	Iframe   string  `json:"iframe" example:"https://example.com/embed/abc123"`
}

type SegmentResponse struct {
	Start int `json:"start" example:"90"`
	End   int `json:"end" example:"180"`
}

type TrackResponse struct {
	URL     string `json:"url" example:"/proxy?p=encodedurl&s=hd"`
	Raw     string `json:"raw" example:"https://example.com/subtitles.vtt"`
	Kind    string `json:"kind" example:"captions"`
	Label   string `json:"label" example:"English"`
	Default bool   `json:"default" example:"true"`
}

type AnimeListResponse = Pagination[AnimeResponse]
type SeasonalAnimeListResponse = []SeasonalAnimeResponse
type EpisodeListResponse = []EpisodeResponse
type EpisodeServerListResponse = []EpisodeServerResponse
type TrendingAnimeListResponse = []AnimeResponse
type PopularAnimeListResponse = []AnimeResponse

