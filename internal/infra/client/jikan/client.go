package jikan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: "https://api.jikan.moe/v4",
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       20,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

func (c *Client) GetAnimeCharacters(ctx context.Context, malID int) (CharactersResponse, error) {
	url := fmt.Sprintf("%s/anime/%d/characters", c.baseURL, malID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jikanResp JikanCharactersResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return nil, err
	}

	var characters CharactersResponse
	for _, char := range jikanResp.Data {
		charNode := CharactersNode{
			MalID:     char.Character.MalID,
			Name:      char.Character.Name,
			Role:      char.Role,
			Favorites: char.Favorites,
			Image:     char.Character.Images.Webp,
		}
		characters = append(characters, charNode)
	}

	return characters, nil
}

func (c *Client) GetCharacterFull(ctx context.Context, malID int) (CharacterFullData, error) {
	url := fmt.Sprintf("%s/characters/%d/full", c.baseURL, malID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return CharacterFullData{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return CharacterFullData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CharacterFullData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jikanResp JikanCharacterFullDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return CharacterFullData{}, err
	}

	return jikanResp.Data, nil
}

func (c *Client) GetPersonFull(ctx context.Context, malID int) (PersonFullData, error) {
	url := fmt.Sprintf("%s/people/%d/full", c.baseURL, malID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return PersonFullData{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return PersonFullData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PersonFullData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jikanResp JikanPersonFullDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&jikanResp); err != nil {
		return PersonFullData{}, err
	}

	return jikanResp.Data, nil
}
