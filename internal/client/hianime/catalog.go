package hianime

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HianimeCatalog struct {
	fetcher *HianimeFetcher
}

func NewCatalog(fetcher *HianimeFetcher) *HianimeCatalog {
	return &HianimeCatalog{
		fetcher: fetcher,
	}
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

func (c *HianimeCatalog) AZList(ctx context.Context, page int) (Pagination[ScrapedAnimeInfoDto], error) {
	headers := map[string]string{
		"Referer":    c.fetcher.baseURL + "/az-list",
		"User-Agent": c.fetcher.randomUA(),
	}
	doc, err := c.fetcher.GetDocument(ctx, "/az-list?page="+strconv.Itoa(page), headers)
	if err != nil {
		return Pagination[ScrapedAnimeInfoDto]{}, err
	}

	return Pagination[ScrapedAnimeInfoDto]{
		PageInfo: extractPageInfo(doc),
		Items:    extractAnimesFromPage(doc),
	}, nil
}

func (c *HianimeCatalog) RecentlyUpdated(ctx context.Context, page int) (Pagination[ScrapedAnimeInfoDto], error) {
	headers := map[string]string{
		"Referer":    c.fetcher.baseURL + "/home",
		"User-Agent": c.fetcher.randomUA(),
	}
	doc, err := c.fetcher.GetDocument(ctx, "/recently-updated?page="+strconv.Itoa(page), headers)
	if err != nil {
		return Pagination[ScrapedAnimeInfoDto]{}, err
	}

	return Pagination[ScrapedAnimeInfoDto]{
		PageInfo: extractPageInfo(doc),
		Items:    extractAnimesFromPage(doc),
	}, nil
}

func (c *HianimeCatalog) AnimeInfo(ctx context.Context, hiAnimeID string) (ScrapedAnimeInfoDto, error) {
	headers := map[string]string{
		"Referer":    c.fetcher.baseURL,
		"User-Agent": c.fetcher.randomUA(),
	}
	doc, err := c.fetcher.GetDocument(ctx, "/"+hiAnimeID, headers)
	if err != nil {
		return ScrapedAnimeInfoDto{}, err
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

	genreSlice := []string{}
	doc.Find(".anisc-info .item-list a").Each(func(_ int, el *goquery.Selection) {
		if href, _ := el.Attr("href"); strings.Contains(href, "genre") {
			genreSlice = append(genreSlice, strings.TrimSpace(el.Text()))
		}
	})

	genre := "Unknown"
	if len(genreSlice) > 0 {
		genre = strings.Join(genreSlice, ", ")
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

func (c *HianimeCatalog) Episodes(ctx context.Context, hiAnimeID string) ([]ScrapedEpisodeDto, error) {
	parts := strings.Split(hiAnimeID, "-")
	eid := parts[len(parts)-1]

	headers := map[string]string{
		"Referer":          c.fetcher.baseURL + "/watch/" + hiAnimeID,
		"User-Agent":       c.fetcher.randomUA(),
		"X-Requested-With": "XMLHttpRequest",
	}

	var wrapper struct {
		HTML string `json:"html"`
	}
	ok, err := c.fetcher.GetAjax(ctx, "/ajax/v2/episode/list/"+eid, headers, &wrapper)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil // No episodes found
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrapper.HTML))
	if err != nil {
		return nil, err
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

func (c *HianimeCatalog) EpisodeServers(ctx context.Context, hiAnimeID, episodeID string) ([]ScrapedEpisodeServerDto, error) {
	headers := map[string]string{
		"Referer":          c.fetcher.baseURL + "/watch/" + hiAnimeID,
		"User-Agent":       c.fetcher.randomUA(),
		"X-Requested-With": "XMLHttpRequest",
	}

	var wrapper struct {
		HTML string `json:"html"`
	}
	ok, err := c.fetcher.GetAjax(ctx, "/ajax/v2/episode/servers?episodeId="+episodeID, headers, &wrapper)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil // No servers found
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrapper.HTML))
	if err != nil {
		return nil, err
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
				if seen["sub"] == struct{}{} {
					continue
				}
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
