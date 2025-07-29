package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type AnilistProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

func NewAnilistProvider(clientID, clientSecret, redirectURL string) *AnilistProvider {
	return &AnilistProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

func (a *AnilistProvider) Name() string {
	return "anilist"
}

func (a *AnilistProvider) AuthURL(ctx context.Context, _ string) (string, error) {
	return fmt.Sprintf(
		"https://anilist.co/api/v2/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		a.clientID,
		url.QueryEscape(a.redirectURL),
	), nil
}

func (a *AnilistProvider) ExchangeToken(ctx context.Context, _ string, code string) (TokenResponse, error) {
	body := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     a.clientID,
		"client_secret": a.clientSecret,
		"redirect_uri":  a.redirectURL,
		"code":          code,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return TokenResponse{}, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://anilist.co/api/v2/oauth/token", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("failed to exchange token")
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return TokenResponse{}, err
	}

	return token, nil
}

var _ Provider = (*AnilistProvider)(nil)
