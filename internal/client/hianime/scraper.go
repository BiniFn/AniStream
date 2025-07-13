package hianime

import (
	"context"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/client/hianime/streaming"
)

type HianimeScraper struct {
	catalog *HianimeCatalog
	stream  *streaming.Fetcher
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
		stream:  streaming.NewFetcher(client),
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

func (s *HianimeScraper) GetStreamingData(
	ctx context.Context,
	serverID, streamType, serverName string,
) (streaming.ScrapedUnencryptedSources, error) {
	return s.stream.Fetch(ctx, serverName, serverID, streamType)
}
