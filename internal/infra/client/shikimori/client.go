package shikimori

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/infra/cache"
)

type Client struct {
	baseURL     string
	httpClient  *http.Client
	redisClient *cache.RedisClient
}

func NewClient(redisClient *cache.RedisClient) *Client {
	return &Client{
		baseURL: "https://shikimori.one/api",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       20,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
		redisClient: redisClient,
	}
}

func (c *Client) GetAnimeFranchise(ctx context.Context, malID int) (*FranchiseResponse, error) {
	var id int
	ok, _ := c.redisClient.Get(ctx, fmt.Sprintf("shikimori:derived_from:%d", malID), &id)
	if ok {
		malID = id
	}

	key := fmt.Sprintf("shikimori:franchise:%d", malID)
	var fr FranchiseResponse
	ok, _ = c.redisClient.Get(ctx, key, &fr)
	if ok {
		return &fr, nil
	}

	url := fmt.Sprintf("%s/animes/%d/franchise", c.baseURL, malID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var franchiseResponse FranchiseResponse
	if err := json.NewDecoder(resp.Body).Decode(&franchiseResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	pipeline := c.redisClient.Pipeline()
	pipeline.Set(ctx, key, franchiseResponse, 7*24*time.Hour)
	for _, node := range franchiseResponse.Nodes {
		pipeline.Set(ctx, fmt.Sprintf("shikimori:derived_from:%d", node.ID), malID, 7*24*time.Hour)
	}
	pipeline.Exec(ctx)

	return &franchiseResponse, nil
}
