package hianime

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type HianimeFetcher struct {
	baseURL string
	client  *http.Client
}

func NewFetcher(baseURL string, client *http.Client) *HianimeFetcher {
	return &HianimeFetcher{
		client:  client,
		baseURL: baseURL,
	}
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
}

func (f *HianimeFetcher) randomUA() string {
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	if err != nil {
		return userAgents[0]
	}

	return userAgents[randomIndex.Int64()]
}

func (f *HianimeFetcher) GetDocument(
	ctx context.Context,
	path string,
	headers map[string]string,
) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

func (f *HianimeFetcher) GetAjax(ctx context.Context, path string, headers map[string]string, dest any) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL+path, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to fetch AJAX content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return false, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return true, nil
}
