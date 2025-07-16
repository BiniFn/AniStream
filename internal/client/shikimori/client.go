package shikimori

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coeeter/aniways/internal/cache"
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
	if ok, err := c.redisClient.Get(ctx, fmt.Sprintf("shikimori:derived_from:%d", malID), &id); err != nil {
		log.Printf("failed to get derived ID from cache: %v", err)
	} else if ok {
		log.Printf("found derived ID in cache for MAL ID %d: %d", malID, id)
		malID = id
	} else {
		log.Printf("cache miss for derived ID with MAL ID %d, using original ID", malID)
	}

	key := fmt.Sprintf("shikimori:franchise:%d", malID)
	var fr FranchiseResponse
	if ok, err := c.redisClient.Get(ctx, key, &fr); err != nil {
		log.Printf("failed to get franchise from cache: %v", err)
	} else if ok {
		log.Printf("found franchise in cache for MAL ID %d", malID)
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
	if _, err := pipeline.Exec(ctx); err != nil {
		log.Printf("failed to execute pipeline for caching franchise: %v", err)
	} else {
		log.Printf("successfully cached franchise for MAL ID %d", malID)
	}

	return &franchiseResponse, nil
}
