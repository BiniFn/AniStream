package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Luzifer/go-openssl/v4"
)

const (
	keyURL = "https://raw.githubusercontent.com/itzzzme/megacloud-keys/refs/heads/main/key.txt"
)

type HianimeProvider struct {
	client *http.Client
}

func NewHianimeProvider(client *http.Client) *HianimeProvider {
	return &HianimeProvider{
		client: client,
	}
}

func (p *HianimeProvider) Supports(serverName string) bool {
	return strings.HasPrefix(strings.ToLower(serverName), "hd")
}

func (p *HianimeProvider) FetchSources(
	ctx context.Context,
	serverID, streamType, serverName string,
) (ScrapedUnencryptedSources, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://hianimez.to/ajax/v2/episode/sources?id=%s", serverID),
		nil,
	)
	req.Header.Set("Referer", "https://hianimez.to/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := p.client.Do(req)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Link string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	if data.Link == "" {
		return ScrapedUnencryptedSources{}, fmt.Errorf("no link found in response")
	}

	parsedURL, err := url.Parse(data.Link)
	if err != nil {
		return ScrapedUnencryptedSources{}, fmt.Errorf("failed to parse streaming link: %w", err)
	}

	parts := strings.Split(parsedURL.Path, "/")
	xrax := parts[len(parts)-1]
	parsedURL.Path = ""
	parsedURL.RawQuery = ""
	origin := parsedURL.String()

	ajaxURL := fmt.Sprintf("%s/embed-2/v2/e-1/getSources?id=%s", origin, xrax)

	ajaxReq, _ := http.NewRequestWithContext(ctx, "GET", ajaxURL, nil)
	ajaxResp, err := p.client.Do(ajaxReq)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	defer ajaxResp.Body.Close()

	var enc ScrapedEncryptedSources
	if err := json.NewDecoder(ajaxResp.Body).Decode(&enc); err != nil {
		return ScrapedUnencryptedSources{}, err
	}

	keyReq, _ := http.NewRequestWithContext(ctx, "GET", keyURL, nil)
	keyResp, err := p.client.Do(keyReq)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	defer keyResp.Body.Close()
	keyBytes, _ := io.ReadAll(keyResp.Body)
	key := strings.TrimSpace(string(keyBytes))

	o := openssl.New()
	plain, err := o.DecryptBytes(
		key,
		[]byte(enc.Sources),
		openssl.BytesToKeyMD5,
	)
	if err != nil {
		return ScrapedUnencryptedSources{}, err
	}

	var list []struct {
		File string `json:"file"`
	}
	if err := json.Unmarshal(plain, &list); err != nil {
		return ScrapedUnencryptedSources{}, err
	}
	if len(list) == 0 {
		return ScrapedUnencryptedSources{}, fmt.Errorf("no sources after decryption")
	}

	return ScrapedUnencryptedSources{
		Source:     list[0].File,
		ServerName: serverName,
		Type:       streamType,
		Intro:      enc.Intro,
		Outro:      enc.Outro,
		Tracks:     enc.Tracks,
	}, nil
}
