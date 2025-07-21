package hianime

import (
	"context"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/client/hianime/streams"
)

type HianimeScraper struct {
	catalog *HianimeCatalog
	stream  *streams.Streams
}

func NewHianimeScraper() *HianimeScraper {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       20,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	return &HianimeScraper{
		catalog: NewCatalog(NewFetcher("https://hianimez.to", client)),
		stream:  streams.NewStreams(client),
	}
}

func (s *HianimeScraper) GetAZList(
	ctx context.Context,
	page int,
) (Pagination[ScrapedAnimeInfoDto], error) {
	return s.catalog.AZList(ctx, page)
}

func (s *HianimeScraper) GetRecentlyUpdatedAnime(
	ctx context.Context,
	page int,
) (Pagination[ScrapedAnimeInfoDto], error) {
	return s.catalog.RecentlyUpdated(ctx, page)
}

func (s *HianimeScraper) GetAnimeInfoByHiAnimeID(
	ctx context.Context,
	hiAnimeID string,
) (ScrapedAnimeInfoDto, error) {
	return s.catalog.AnimeInfo(ctx, hiAnimeID)
}

func (s *HianimeScraper) GetAnimeEpisodes(
	ctx context.Context,
	hiAnimeID string,
) ([]ScrapedEpisodeDto, error) {
	return s.catalog.Episodes(ctx, hiAnimeID)
}

func (s *HianimeScraper) GetEpisodeServers(
	ctx context.Context,
	hiAnimeID, episodeID string,
) ([]ScrapedEpisodeServerDto, error) {
	return s.catalog.EpisodeServers(ctx, hiAnimeID, episodeID)
}

func (s *HianimeScraper) GetEpisodeLangs(
	ctx context.Context,
	hiAnimeID, episodeID string,
) ([]string, error) {
	return s.catalog.EpisodeLangs(ctx, hiAnimeID, episodeID)
}

func (s *HianimeScraper) GetEpisodeStream(
	ctx context.Context,
	episodeID, streamType string,
) (string, error) {
	return s.stream.GetStreamingSource(ctx, episodeID, streamType)
}
