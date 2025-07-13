package anilist

import (
	"context"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	operations "github.com/coeeter/aniways/internal/client/anilist/graphql"
)

var (
	url = "https://graphql.anilist.co"
)

type Client struct {
	graphqlClient graphql.Client
}

func New() *Client {
	return &Client{
		graphqlClient: graphql.NewClient(url, http.DefaultClient),
	}
}

func (c *Client) GetSeasonalMedia(ctx context.Context, year int, season string) (operations.GetSeasonalAnimeResponse, error) {
	data, err := operations.GetSeasonalAnime(ctx, c.graphqlClient, year, operations.MediaSeason(season))
	if err != nil {
		return operations.GetSeasonalAnimeResponse{}, err
	}
	return operations.GetSeasonalAnimeResponse{
		Page: data.Page,
	}, nil
}

func (c *Client) GetTrendingAnime(ctx context.Context) (operations.GetTrendingAnimeResponse, error) {
	data, err := operations.GetTrendingAnime(ctx, c.graphqlClient)
	if err != nil {
		return operations.GetTrendingAnimeResponse{}, err
	}
	return operations.GetTrendingAnimeResponse{
		Page: data.Page,
	}, nil
}

func (c *Client) GetPopularAnime(ctx context.Context) (operations.GetPopularAnimeResponse, error) {
	data, err := operations.GetPopularAnime(ctx, c.graphqlClient)
	if err != nil {
		return operations.GetPopularAnimeResponse{}, err
	}
	return operations.GetPopularAnimeResponse{
		Page: data.Page,
	}, nil
}

func (c *Client) GetAnimeDetails(ctx context.Context, id int) (operations.GetAnimeDetailsResponse, error) {
	data, err := operations.GetAnimeDetails(ctx, c.graphqlClient, id)
	if err != nil {
		return operations.GetAnimeDetailsResponse{}, err
	}
	return operations.GetAnimeDetailsResponse{
		Media: data.Media,
	}, nil
}
