package myanimelist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client interface {
	GetAnimeMetadata(ctx context.Context, malID int) (*MalAnimeMetadata, error)
	GetTrailer(ctx context.Context, malID int) (string, error)
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
}

func NewClient(config ClientConfig) Client {
	return &client{
		baseURL:  "https://api.myanimelist.net/v2",
		ClientID: config.ClientID,
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

type client struct {
	baseURL    string
	ClientID   string
	httpClient *http.Client
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

func (c *client) GetAnimeMetadata(ctx context.Context, malID int) (*MalAnimeMetadata, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/anime/%d", c.baseURL, malID))
	q := u.Query()
	q.Set("fields", strings.Join(metadataFields, ","))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-MAL-CLIENT-ID", c.ClientID)

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

func (c *client) GetTrailer(ctx context.Context, malID int) (string, error) {
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
	if !strings.Contains(trailerURL, "youtube.com") {
		return "", fmt.Errorf("trailer is not a YouTube link: %s", trailerURL)
	}
	return trailerURL, nil
}
