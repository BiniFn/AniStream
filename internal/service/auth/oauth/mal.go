package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coeeter/aniways/internal/cache"
	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type MALProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	repo         *repository.Queries
	redis        *cache.RedisClient
}

func NewMALProvider(clientID, clientSecret, redirectURL string, repo *repository.Queries, redis *cache.RedisClient) *MALProvider {
	return &MALProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		repo:         repo,
		redis:        redis,
	}
}

func (m *MALProvider) Name() string {
	return "myanimelist"
}

func (m *MALProvider) AuthURL(ctx context.Context, state string) (string, error) {
	verifier, _ := generateCodeVerifier(generateCodeVerifierParams{})

	key := fmt.Sprintf("oauth:mal:%s", state)
	if err := m.redis.Set(ctx, key, verifier, 10*time.Minute); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&state=%s&code_challenge=%s&code_challenge_method=plain&redirect_uri=%s",
		url.QueryEscape(m.clientID),
		url.QueryEscape(state),
		url.QueryEscape(verifier),
		url.QueryEscape(m.redirectURL),
	), nil
}

func (m *MALProvider) ExchangeToken(ctx context.Context, userID, state, code string) error {
	key := fmt.Sprintf("oauth:mal:%s", state)

	var verifier string
	ok, err := m.redis.Get(ctx, key, &verifier)
	if err != nil {
		return fmt.Errorf("redis get error: %w", err)
	}
	if !ok {
		return fmt.Errorf("invalid or expired state")
	}

	_ = m.redis.Del(ctx, key)

	form := url.Values{}
	form.Add("client_id", m.clientID)
	form.Add("client_secret", m.clientSecret)
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", m.redirectURL)
	form.Add("code_verifier", verifier)

	req, _ := http.NewRequest("POST", "https://myanimelist.net/v1/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tokenResponse := TokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	expiresAt := pgtype.Timestamp{
		Time:  time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		Valid: true,
	}

	return m.repo.SaveOauthToken(ctx, repository.SaveOauthTokenParams{
		UserID:       userID,
		Token:        tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		Provider:     repository.Provider(m.Name()),
		ExpiresAt:    expiresAt,
	})
}

func (m *MALProvider) RefreshToken(ctx context.Context, userID, refreshToken string) error {
	form := url.Values{}
	form.Add("client_id", m.clientID)
	form.Add("client_secret", m.clientSecret)
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", "https://myanimelist.net/v1/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tokenResponse := TokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	expiresAt := pgtype.Timestamp{
		Time:  time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
		Valid: true,
	}

	_, err = m.repo.GetToken(ctx, repository.GetTokenParams{
		UserID:   userID,
		Provider: repository.Provider(m.Name()),
	})

	if err == nil {
		return m.repo.UpdateOauthToken(ctx, repository.UpdateOauthTokenParams{
			UserID:       userID,
			Token:        tokenResponse.AccessToken,
			RefreshToken: tokenResponse.RefreshToken,
			Provider:     repository.Provider(m.Name()),
			ExpiresAt:    expiresAt,
		})
	}

	return m.repo.UpdateOauthToken(ctx, repository.UpdateOauthTokenParams{
		UserID:       userID,
		Token:        tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		Provider:     repository.Provider(m.Name()),
		ExpiresAt:    expiresAt,
	})
}

var _ Provider = (*MALProvider)(nil)
