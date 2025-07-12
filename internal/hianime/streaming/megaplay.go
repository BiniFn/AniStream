package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MegaplayProvider struct {
	client *http.Client
}

func NewMegaplayProvider(client *http.Client) *MegaplayProvider {
	return &MegaplayProvider{client: client}
}

func (p *MegaplayProvider) Supports(serverName string) bool {
	return strings.Contains(strings.ToLower(serverName), "megaplay")
}

func (p *MegaplayProvider) FetchSources(
	ctx context.Context,
	serverID, streamType, _ string,
) (ScrapedUnencryptedSources, error) {
	url := fmt.Sprintf("https://megaplay.buzz/stream/s-2/%s/%s", serverID, streamType)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Referer", "https://megaplay.buzz")
	resp, err := p.client.Do(req)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	mediaID, exists := doc.Find("#megaplay-player").Attr("data-id")
	if mediaID == "" || !exists {
		return ScrapedUnencryptedSources{}, fmt.Errorf("no media ID found in Megaplay streaming data")
	}

	ajax := fmt.Sprintf("https://megaplay.buzz/stream/getSources?id=%s", mediaID)
	ajaxReq, _ := http.NewRequestWithContext(ctx, "GET", ajax, nil)
	ajaxReq.Header.Set("Referer", "https://megaplay.buzz")
	ajaxReq.Header.Set("Origin", url)
	ajaxReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	ajaxResp, err := p.client.Do(ajaxReq)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	defer ajaxResp.Body.Close()

	var meta ScrapedMegaplaySources
	if err := json.NewDecoder(ajaxResp.Body).Decode(&meta); err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	if meta.Sources.File == "" {
		return ScrapedUnencryptedSources{}, fmt.Errorf("no sources in Megaplay response")
	}

	return ScrapedUnencryptedSources{
		Source:     meta.Sources.File,
		ServerName: "Megaplay",
		Type:       streamType,
		Intro:      meta.Intro,
		Outro:      meta.Outro,
		Tracks:     meta.Tracks,
	}, nil
}
