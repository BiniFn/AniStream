package streams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Streams struct {
	client *http.Client
}

func NewStreams(client *http.Client) *Streams {
	return &Streams{client: client}
}

func (s *Streams) GetStreamingSource(ctx context.Context, episodeID, streamType string) (string, error) {
	url := fmt.Sprintf("https://megaplay.buzz/stream/s-2/%s/%s", episodeID, streamType)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Referer", "https://megaplay.buzz")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	mediaID, exists := doc.Find("#megaplay-player").Attr("data-id")
	if mediaID == "" || !exists {
		return "", fmt.Errorf("no media ID found in Megaplay streaming data")
	}

	ajax := fmt.Sprintf("https://megaplay.buzz/stream/getSources?id=%s", mediaID)
	ajaxReq, _ := http.NewRequestWithContext(ctx, "GET", ajax, nil)
	ajaxReq.Header.Set("Referer", "https://megaplay.buzz")
	ajaxReq.Header.Set("Origin", url)
	ajaxReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	ajaxResp, err := s.client.Do(ajaxReq)
	if err != nil {
		return "", err
	}
	defer ajaxResp.Body.Close()

	var meta struct {
		Sources struct {
			File string `json:"file"`
		} `json:"sources"`
	}
	if err := json.NewDecoder(ajaxResp.Body).Decode(&meta); err != nil {
		return "", err
	}
	if meta.Sources.File == "" {
		return "", fmt.Errorf("no sources in Megaplay response")
	}

	return meta.Sources.File, nil
}
