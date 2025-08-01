package anilist

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"
	operations "github.com/coeeter/aniways/internal/infra/client/anilist/graphql"
	"github.com/coeeter/aniways/internal/repository"
)

var (
	url = "https://graphql.anilist.co"
)

type Client struct {
	graphqlClient graphql.Client
}

type httpClient struct {
	client *http.Client
}

func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	if token, ok := req.Context().Value("token").(string); ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return h.client.Do(req)
}

func New() *Client {
	return &Client{
		graphqlClient: graphql.NewClient(url, &httpClient{
			client: http.DefaultClient,
		}),
	}
}

func (c *Client) GetSeasonalMedia(ctx context.Context, year int, season string) (operations.GetSeasonalAnimeResponse, error) {
	data, err := operations.GetSeasonalAnime(ctx, c.graphqlClient, year, operations.MediaSeason(season))
	if err != nil {
		return operations.GetSeasonalAnimeResponse{}, err
	}
	return *data, nil
}

func (c *Client) GetTrendingAnime(ctx context.Context) (operations.GetTrendingAnimeResponse, error) {
	data, err := operations.GetTrendingAnime(ctx, c.graphqlClient)
	if err != nil {
		return operations.GetTrendingAnimeResponse{}, err
	}
	return *data, nil
}

func (c *Client) GetPopularAnime(ctx context.Context) (operations.GetPopularAnimeResponse, error) {
	data, err := operations.GetPopularAnime(ctx, c.graphqlClient)
	if err != nil {
		return operations.GetPopularAnimeResponse{}, err
	}
	return *data, nil
}

func (c *Client) GetAnimeDetails(ctx context.Context, id int) (operations.GetAnimeDetailsResponse, error) {
	data, err := operations.GetAnimeDetails(ctx, c.graphqlClient, id)
	if err != nil {
		return operations.GetAnimeDetailsResponse{}, err
	}
	return *data, nil
}

type GetUserAnimeListParams struct {
	Token        string
	Page         int
	ItemsPerPage int
}

var ErrInvalidToken = errors.New("invalid token")

func (c *Client) extractUserID(token string) (int, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, ErrInvalidToken
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, ErrInvalidToken
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int(userID), nil
}

func (c *Client) withToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, "token", token)
}

func (c *Client) GetUserAnimeList(ctx context.Context, params GetUserAnimeListParams) (operations.GetUserAnimeListResponse, error) {
	userID, err := c.extractUserID(params.Token)
	if err != nil {
		return operations.GetUserAnimeListResponse{}, err
	}

	ctx = c.withToken(ctx, params.Token)

	data, err := operations.GetUserAnimeList(ctx, c.graphqlClient, params.Page, params.ItemsPerPage, userID)
	if err != nil {
		return operations.GetUserAnimeListResponse{}, err
	}

	return *data, nil
}

type UpdateAnimeListParams struct {
	Token           string
	MalID           int
	AnilistID       int
	Status          string
	WatchedEpisodes int
}

func (c *Client) convertFromRepoStatus(status string) operations.MediaListStatus {
	switch status {
	case string(repository.LibraryStatusWatching):
		return operations.MediaListStatusCurrent
	case string(repository.LibraryStatusCompleted):
		return operations.MediaListStatusCompleted
	case string(repository.LibraryStatusPaused):
		return operations.MediaListStatusPaused
	case string(repository.LibraryStatusDropped):
		return operations.MediaListStatusDropped
	case string(repository.LibraryStatusPlanning):
		return operations.MediaListStatusPlanning
	default:
		return operations.MediaListStatusPlanning
	}
}

func (c *Client) UpdateAnimeList(ctx context.Context, params UpdateAnimeListParams) error {
	ctx = c.withToken(ctx, params.Token)

	status := c.convertFromRepoStatus(params.Status)
	watchedEpisodes := params.WatchedEpisodes

	res, err := operations.GetAnimeId(ctx, c.graphqlClient, []int{params.AnilistID}, []int{params.MalID})
	if err != nil {
		return err
	}

	_, err = operations.SaveMediaList(ctx, c.graphqlClient, res.Media.GetId(), status, watchedEpisodes)
	if err != nil {
		return err
	}

	return nil
}

type DeleteAnimeListParams struct {
	Token     string
	MalID     int
	AnilistID int
}

func (c *Client) DeleteAnimeList(ctx context.Context, params DeleteAnimeListParams) error {
	ctx = c.withToken(ctx, params.Token)

	userId, err := c.extractUserID(params.Token)
	if err != nil {
		return err
	}

	res, err := operations.GetAnimeId(ctx, c.graphqlClient, []int{params.AnilistID}, []int{params.MalID})
	if err != nil {
		return err
	}

	entry, err := operations.GetUserEntryId(ctx, c.graphqlClient, res.Media.GetId(), userId)
	if err != nil {
		return err
	}

	_, err = operations.DeleteMediaListEntry(ctx, c.graphqlClient, entry.MediaList.GetId())
	if err != nil {
		return err
	}

	return nil
}
