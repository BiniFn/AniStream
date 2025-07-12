package hianime

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HianimeScraper struct {
	client  *http.Client
	baseURL string
}

func NewHianimeScraper() *HianimeScraper {
	return &HianimeScraper{
		baseURL: "https://hianimez.to",
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

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
}

func (s *HianimeScraper) randomUA() string {

	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	if err != nil {
		return userAgents[0]
	}

	return userAgents[randomIndex.Int64()]
}

func (s *HianimeScraper) fetchHTML(ctx context.Context, path, refererPath string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Referer", s.baseURL+refererPath)
	req.Header.Set("User-Agent", s.randomUA())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *HianimeScraper) fetchAjaxHTML(ctx context.Context, path, refererPath string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Referer", s.baseURL+refererPath)
	req.Header.Set("User-Agent", s.randomUA())
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AJAX content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var wrapper struct {
		HTML string `json:"html"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode AJAX response: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrapper.HTML))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML from AJAX response: %w", err)
	}
	return doc, nil
}

func extractPageInfo(d *goquery.Document) PageInfo {
	p := d.Find(".pagination")
	if p.Length() == 0 {
		return PageInfo{TotalPages: 1, CurrentPage: 1, HasNextPage: false, HasPreviousPage: false}
	}

	li := p.Find("li")
	lastLi := li.Last()

	var lastPage int
	if lastLi.HasClass("active") {
		lastPage, _ = strconv.Atoi(strings.TrimSpace(lastLi.Text()))
		if lastPage < 1 {
			lastPage = 1
		}
	} else {
		href, ok := lastLi.Find("a").Attr("href")
		if ok {
			parts := strings.Split(href, "page=")
			if n, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				lastPage = n
			}
		}
		if lastPage < 1 {
			lastPage = 1
		}
	}

	currPage := 1
	if curText := strings.TrimSpace(p.Find("li.active a").Text()); curText != "" {
		if n, err := strconv.Atoi(curText); err == nil {
			currPage = n
		}
	}

	return PageInfo{
		TotalPages:      lastPage,
		CurrentPage:     currPage,
		HasNextPage:     currPage < lastPage,
		HasPreviousPage: currPage > 1,
	}
}

func extractAnimesFromPage(d *goquery.Document) []ScrapedAnimeInfoDto {
	out := []ScrapedAnimeInfoDto{}
	d.Find("div.flw-item").Each(func(_ int, el *goquery.Selection) {
		href, _ := el.Find(".film-poster a").Attr("href")
		parts := strings.Split(strings.Trim(href, "/"), "/")
		id := parts[len(parts)-1]

		ename := strings.TrimSpace(el.Find(".film-detail .film-name a").Text())
		jname, _ := el.Find(".film-detail .film-name a").Attr("data-jname")
		poster, _ := el.Find(".film-poster img").Attr("data-src")
		episodesStr := strings.TrimSpace(el.Find(".film-poster .tick-sub").Text())
		episodes, _ := strconv.Atoi(episodesStr)

		out = append(out, ScrapedAnimeInfoDto{
			HiAnimeID:   id,
			EName:       ename,
			JName:       jname,
			PosterURL:   poster,
			LastEpisode: episodes,
		})
	})
	return out
}

func (s *HianimeScraper) GetAZList(ctx context.Context, page int) (Pagination[ScrapedAnimeInfoDto], error) {
	doc, err := s.fetchHTML(ctx, fmt.Sprintf("/az-list?page=%d", page), "/az-list")
	if err != nil {
		return Pagination[ScrapedAnimeInfoDto]{}, fmt.Errorf("failed to fetch AZ list: %w", err)
	}
	return Pagination[ScrapedAnimeInfoDto]{
		PageInfo: extractPageInfo(doc),
		Items:    extractAnimesFromPage(doc),
	}, nil
}

func (s *HianimeScraper) GetRecentlyUpdatedAnime(ctx context.Context, page int) (Pagination[ScrapedAnimeInfoDto], error) {
	doc, err := s.fetchHTML(ctx, fmt.Sprintf("/recently-updated?page=%d", page), "/home")
	if err != nil {
		return Pagination[ScrapedAnimeInfoDto]{}, fmt.Errorf("failed to fetch recently updated anime: %w", err)
	}
	return Pagination[ScrapedAnimeInfoDto]{
		PageInfo: extractPageInfo(doc),
		Items:    extractAnimesFromPage(doc),
	}, nil
}

func (s *HianimeScraper) GetAnimeInfoByHiAnimeID(ctx context.Context, hiAnimeID string) (ScrapedAnimeInfoDto, error) {
	doc, err := s.fetchHTML(ctx, "/"+hiAnimeID, "")
	if err != nil {
		return ScrapedAnimeInfoDto{}, fmt.Errorf("failed to fetch anime info: %w", err)
	}

	syncJSON := doc.Find("#syncData").Text()
	var sync struct {
		MalID     string `json:"mal_id"`
		AnilistID string `json:"anilist_id"`
	}
	_ = json.Unmarshal([]byte(syncJSON), &sync)
	malID, _ := strconv.Atoi(sync.MalID)
	anilistID, _ := strconv.Atoi(sync.AnilistID)

	titleEl := doc.Find("h2.film-name.dynamic-name")
	ename := strings.TrimSpace(titleEl.Text())
	jname, _ := titleEl.Attr("data-jname")
	poster, _ := doc.Find(".film-poster img").Attr("src")

	genre_slice := []string{}
	doc.Find(".anisc-info .item-list a").Each(func(_ int, el *goquery.Selection) {
		if href, _ := el.Attr("href"); strings.Contains(href, "genre") {
			genre_slice = append(genre_slice, strings.TrimSpace(el.Text()))
		}
	})

	genre := "Unknown"
	if len(genre_slice) > 0 {
		genre = strings.Join(genre_slice, ", ")
	}

	lastEpTxt := doc.Find(".tick-item.tick-sub").First().Text()
	lastEp, _ := strconv.Atoi(strings.TrimSpace(lastEpTxt))

	return ScrapedAnimeInfoDto{
		HiAnimeID:   hiAnimeID,
		EName:       ename,
		JName:       jname,
		PosterURL:   poster,
		Genre:       genre,
		MalID:       malID,
		AnilistID:   anilistID,
		LastEpisode: lastEp,
	}, nil
}

func (s *HianimeScraper) GetAnimeEpisodes(ctx context.Context, hiAnimeID string) ([]ScrapedEpisodeDto, error) {
	parts := strings.Split(hiAnimeID, "-")
	path := fmt.Sprintf("/ajax/v2/episode/list/%s", parts[len(parts)-1])
	referer := fmt.Sprintf("/watch/%s", hiAnimeID)
	doc, err := s.fetchAjaxHTML(ctx, path, referer)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch episodes: %w", err)
	}

	var eps []ScrapedEpisodeDto
	doc.Find(".detail-infor-content .ss-list a").Each(func(_ int, el *goquery.Selection) {
		title := strings.TrimSpace(el.AttrOr("title", ""))
		num, _ := strconv.Atoi(el.AttrOr("data-number", "1"))
		filler := el.HasClass("ssl-item-filler")
		href, _ := el.Attr("href")
		epID := strings.TrimSpace(strings.Split(href, "?ep=")[1])
		eps = append(eps, ScrapedEpisodeDto{
			EpisodeID: epID,
			Title:     title,
			Number:    num,
			IsFiller:  filler,
		})
	})

	return eps, nil
}

func (s *HianimeScraper) GetEpisodeServers(ctx context.Context, hiAnimeID, episodeID string) ([]ScrapedEpisodeServerDto, error) {
	path := fmt.Sprintf("/ajax/v2/episode/servers?episodeId=%s", episodeID)
	referer := fmt.Sprintf("/watch/%s", hiAnimeID)
	doc, err := s.fetchAjaxHTML(ctx, path, referer)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch episode servers: %w", err)
	}

	allowedIDs := map[string]struct{}{
		"1": {}, "4": {}, "6": {},
	}

	var out []ScrapedEpisodeServerDto
	seen := make(map[string]struct{})

	doc.Find(".server-item").Each(func(_ int, el *goquery.Selection) {
		sid := el.AttrOr("data-server-id", "")
		if _, ok := allowedIDs[sid]; !ok {
			return
		}

		t := el.AttrOr("data-type", "")
		name := strings.TrimSpace(el.Text())
		id := el.AttrOr("data-id", "")

		out = append(out, ScrapedEpisodeServerDto{
			Type:       t,
			ServerName: name,
			ServerID:   id,
		})
		seen[t] = struct{}{}
	})

	fallbackTypes := []string{"sub", "dub", "raw"}
	for _, t := range fallbackTypes {
		if _, ok := seen[t]; ok {
			if t == "raw" {
				t = "sub"
			}

			out = append(out, ScrapedEpisodeServerDto{
				Type:       t,
				ServerName: "Megaplay",
				ServerID:   episodeID,
			})
		}
	}

	return out, nil
}
