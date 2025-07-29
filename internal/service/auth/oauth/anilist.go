package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/coeeter/aniways/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type AnilistProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	repo         *repository.Queries
}

func NewAnilistProvider(clientID, clientSecret, redirectURL string, repo *repository.Queries) *AnilistProvider {
	return &AnilistProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		repo:         repo,
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

func (a *AnilistProvider) ExchangeToken(ctx context.Context, userID, _, code string) error {
	body := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     a.clientID,
		"client_secret": a.clientSecret,
		"redirect_uri":  a.redirectURL,
		"code":          code,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://anilist.co/api/v2/oauth/token", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to exchange token")
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return err
	}

	expiresAt := pgtype.Timestamp{
		Time:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		Valid: true,
	}

	_, err = a.repo.GetToken(ctx, repository.GetTokenParams{
		UserID:   userID,
		Provider: repository.Provider(a.Name()),
	})

	if err == nil {
		return a.repo.UpdateOauthToken(ctx, repository.UpdateOauthTokenParams{
			UserID:       userID,
			Token:        token.AccessToken,
			RefreshToken: token.RefreshToken,
			Provider:     repository.Provider(a.Name()),
			ExpiresAt:    expiresAt,
		})
	}

	return a.repo.SaveOauthToken(ctx, repository.SaveOauthTokenParams{
		UserID:       userID,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
		Provider:     repository.Provider(a.Name()),
		ExpiresAt:    expiresAt,
	})
}

var ErrUnsupportedOperation = errors.New("unsupported operation")

func (a *AnilistProvider) RefreshToken(ctx context.Context, userID, refreshToken string) error {
	// anilist does not support refresh tokens
	return ErrUnsupportedOperation
}

var _ Provider = (*AnilistProvider)(nil)
