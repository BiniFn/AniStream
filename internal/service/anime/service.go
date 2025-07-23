package anime

import (
	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/client/anilist"
	"github.com/coeeter/aniways/internal/client/hianime"
	"github.com/coeeter/aniways/internal/client/myanimelist"
	"github.com/coeeter/aniways/internal/client/shikimori"
	"github.com/coeeter/aniways/internal/repository"
)

type AnimeService struct {
	repo            *repository.Queries
	refresher       *MetadataRefresher
	scraper         *hianime.HianimeScraper
	malClient       *myanimelist.Client
	anilistClient   *anilist.Client
	shikimoriClient *shikimori.Client
	redis           *cache.RedisClient
}

func NewAnimeService(
	repo *repository.Queries,
	refresher *MetadataRefresher,
	malClient *myanimelist.Client,
	anilistClient *anilist.Client,
	shikimoriClient *shikimori.Client,
	redis *cache.RedisClient,
) *AnimeService {
	return &AnimeService{
		repo:            repo,
		refresher:       refresher,
		malClient:       malClient,
		anilistClient:   anilistClient,
		shikimoriClient: shikimoriClient,
		scraper:         hianime.NewHianimeScraper(),
		redis:           redis,
	}
}
