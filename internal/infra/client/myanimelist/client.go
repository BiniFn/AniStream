package myanimelist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	baseURL    string
	clientId   string
	httpClient *http.Client
}

func NewClient(clientID string) *Client {
	return &Client{
		baseURL:  "https://api.myanimelist.net/v2",
		clientId: clientID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       20,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

var metadataFields = []string{
	"alternative_titles",
	"synopsis",
	"main_picture",
	"media_type",
	"rating",
	"average_episode_duration",
	"status",
	"num_episodes",
	"studios",
	"rank",
	"mean",
	"num_scoring_users",
	"popularity",
	"start_date",
	"end_date",
	"source",
	"start_season",
}

func (c *Client) GetAnimeMetadata(ctx context.Context, malID int) (*MalAnimeMetadata, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/anime/%d", c.baseURL, malID))
	q := u.Query()
	q.Set("fields", strings.Join(metadataFields, ","))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-MAL-CLIENT-ID", c.clientId)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metadata MalAnimeMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if metadata.MalID != malID {
		return nil, fmt.Errorf("mal ID mismatch: expected %d, got %d", malID, metadata.MalID)
	}

	return &metadata, nil
}

func (c *Client) GetTrailer(ctx context.Context, malID int) (string, error) {
	url := fmt.Sprintf("https://myanimelist.net/anime/%d", malID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	trailerURL, exists := document.Find("a.iframe").Attr("href")
	if !exists {
		return "", fmt.Errorf("trailer not found for mal ID %d", malID)
	}

	return trailerURL, nil
}

type GetAnimeListParams struct {
	Token        string
	Status       string
	Sort         string
	Page         int
	ItemsPerPage int
}

func (c *Client) GetAnimeList(ctx context.Context, params GetAnimeListParams) (AnimeList, error) {
	u, err := url.Parse(fmt.Sprintf("%s/users/@me/animelist", c.baseURL))
	if err != nil {
		return AnimeList{}, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := url.Values{}

	if params.Status != "" {
		query.Set("status", params.Status)
	}

	if params.Sort != "" {
		query.Set("sort", params.Sort)
	}

	limit := params.ItemsPerPage
	if limit == 0 {
		limit = 30
	}

	offset := (params.Page - 1) * limit

	query.Set("limit", strconv.Itoa(limit))
	query.Set("offset", strconv.Itoa(offset))
	query.Set("fields", "list_status")

	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return AnimeList{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+params.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AnimeList{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return AnimeList{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var animeList AnimeList
	if err := json.NewDecoder(resp.Body).Decode(&animeList); err != nil {
		return AnimeList{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return animeList, nil
}

type UpdateAnimeListParams struct {
	Token           string
	AnimeID         int
	Status          string
	WatchedEpisodes int
}

func (c *Client) UpdateAnimeList(ctx context.Context, params UpdateAnimeListParams) error {
	body := url.Values{}

	s := MalListStatus("").FromRepository(params.Status)
	if s.IsValid() && params.Status != "" {
		body.Set("status", string(s))
	}
	if params.WatchedEpisodes >= 0 {
		body.Set("num_watched_episodes", strconv.Itoa(params.WatchedEpisodes))
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fmt.Sprintf("%s/anime/%d/my_list_status", c.baseURL, params.AnimeID), strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+params.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

type UpdateAnimeListStatusParams struct {
	Token   string
	AnimeID int
	Status  string
}

func (c *Client) UpdateAnimeListStatus(ctx context.Context, params UpdateAnimeListStatusParams) error {
	body := url.Values{}
	s := MalListStatus("").FromRepository(params.Status)
	if s.IsValid() && params.Status != "" {
		body.Set("status", string(s))
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fmt.Sprintf("%s/anime/%d/my_list_status", c.baseURL, params.AnimeID), strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+params.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

type DeleteAnimeListParams struct {
	Token   string
	AnimeID int
}

func (c *Client) DeleteAnimeList(ctx context.Context, params DeleteAnimeListParams) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/anime/%d/my_list_status", c.baseURL, params.AnimeID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+params.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
