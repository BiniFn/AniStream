package anilist

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/Khan/genqlient/graphql"
	operations "github.com/coeeter/aniways/internal/infra/client/anilist/graphql"
	"github.com/coeeter/aniways/internal/repository"
)

const apiURL = "https://graphql.anilist.co"

var ErrInvalidToken = errors.New("invalid token")

type Client struct {
	graphqlClient graphql.Client

	mu         sync.RWMutex
	tokenCache map[string]int
	malIDCache map[int]int
	entryCache map[int]int
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
		graphqlClient: graphql.NewClient(apiURL, &httpClient{client: http.DefaultClient}),
		tokenCache:    make(map[string]int),
		malIDCache:    make(map[int]int),
		entryCache:    make(map[int]int),
	}
}

func (c *Client) withToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, "token", token)
}

func (c *Client) extractUserID(token string) (int, error) {
	c.mu.RLock()
	if id, ok := c.tokenCache[token]; ok {
		c.mu.RUnlock()
		return id, nil
	}
	c.mu.RUnlock()

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

	c.mu.Lock()
	c.tokenCache[token] = int(userID)
	c.mu.Unlock()

	return int(userID), nil
}

func (c *Client) getAnilistIDFromMALID(ctx context.Context, malID int) (int, error) {
	c.mu.RLock()
	if id, ok := c.malIDCache[malID]; ok {
		c.mu.RUnlock()
		return id, nil
	}
	c.mu.RUnlock()

	res, err := operations.GetAnimeId(ctx, c.graphqlClient, malID)
	if err != nil {
		return 0, err
	}

	id := res.Media.GetId()

	c.mu.Lock()
	c.malIDCache[malID] = id
	c.mu.Unlock()

	return id, nil
}

func (c *Client) setEntryCache(malID, entryID int) {
	c.mu.Lock()
	c.entryCache[malID] = entryID
	c.mu.Unlock()
}

func (c *Client) getEntryCache(malID int) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.entryCache[malID]
	return id, ok
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

type InsertAnimeToListParams struct {
	Token           string
	MalID           int
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
	default:
		return operations.MediaListStatusPlanning
	}
}

func (c *Client) InsertAnimeToList(ctx context.Context, params InsertAnimeToListParams) error {
	ctx = c.withToken(ctx, params.Token)

	mediaID, err := c.getAnilistIDFromMALID(ctx, params.MalID)
	if err != nil {
		return err
	}

	status := c.convertFromRepoStatus(params.Status)

	res, err := operations.InsertMediaListEntry(ctx, c.graphqlClient, mediaID, status, params.WatchedEpisodes)
	if err != nil {
		return err
	}

	c.setEntryCache(params.MalID, res.SaveMediaListEntry.GetId())
	return nil
}

type UpdateAnimeEntryStatusParams struct {
	Token  string
	MalID  int
	Status string
}

func (c *Client) UpdateAnimeEntryStatus(ctx context.Context, params UpdateAnimeEntryStatusParams) error {
	ctx = c.withToken(ctx, params.Token)

	mediaID, err := c.getAnilistIDFromMALID(ctx, params.MalID)
	if err != nil {
		return err
	}

	status := c.convertFromRepoStatus(params.Status)

	res, err := operations.UpdateMediaListStatus(ctx, c.graphqlClient, mediaID, status)
	if err != nil {
		return err
	}

	c.setEntryCache(params.MalID, res.SaveMediaListEntry.GetId())
	return nil
}

type UpdateAnimeEntryProgressParams struct {
	Token           string
	MalID           int
	WatchedEpisodes int
}

func (c *Client) UpdateAnimeEntryProgress(ctx context.Context, params UpdateAnimeEntryProgressParams) error {
	ctx = c.withToken(ctx, params.Token)

	mediaID, err := c.getAnilistIDFromMALID(ctx, params.MalID)
	if err != nil {
		return err
	}

	res, err := operations.UpdateMediaListProgress(ctx, c.graphqlClient, mediaID, params.WatchedEpisodes)
	if err != nil {
		return err
	}

	c.setEntryCache(params.MalID, res.SaveMediaListEntry.GetId())
	return nil
}

type DeleteAnimeListParams struct {
	Token string
	MalID int
}

func (c *Client) DeleteAnimeList(ctx context.Context, params DeleteAnimeListParams) error {
	ctx = c.withToken(ctx, params.Token)

	userID, err := c.extractUserID(params.Token)
	if err != nil {
		return err
	}

	mediaID, err := c.getAnilistIDFromMALID(ctx, params.MalID)
	if err != nil {
		return err
	}

	entryID, ok := c.getEntryCache(params.MalID)
	if !ok {
		entry, err := operations.GetUserEntryId(ctx, c.graphqlClient, mediaID, userID)
		if err != nil {
			return err
		}
		entryID = entry.MediaList.GetId()
		c.setEntryCache(params.MalID, entryID)
	}

	_, err = operations.DeleteMediaListEntry(ctx, c.graphqlClient, entryID)
	return err
}
