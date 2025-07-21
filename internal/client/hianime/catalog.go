package hianime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
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

func (c *HianimeCatalog) EpisodeLangs(ctx context.Context, hiAnimeID, episodeID string) ([]string, error) {
	servers, err := c.EpisodeServers(ctx, hiAnimeID, episodeID)
	if err != nil {
		return nil, err
	}

	var langs []string
	for _, server := range servers {
		switch server.Type {
		case "sub", "raw":
			if exists := slices.Contains(langs, "sub"); exists {
				continue
			}
			langs = append(langs, "sub")
		case "dub":
			if exists := slices.Contains(langs, "dub"); exists {
				continue
			}
			langs = append(langs, "dub")
		default:
			if exists := slices.Contains(langs, "unknown"); exists {
				continue
			}
			langs = append(langs, "unknown")
		}
	}

	return langs, nil
}

func (c *HianimeCatalog) StreamSource(ctx context.Context, episodeID, streamType string) (string, error) {
	url := fmt.Sprintf("https://megaplay.buzz/stream/s-2/%s/%s", episodeID, streamType)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Referer", "https://megaplay.buzz")
	resp, err := c.fetcher.Client.Do(req)
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
	ajaxResp, err := c.fetcher.Client.Do(ajaxReq)
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

func (c *HianimeCatalog) StreamMetadata(
	ctx context.Context,
	hianimeID, episodeID, streamType string,
) (ScrapedStreamMetadata, error) {
	servers, err := c.EpisodeServers(ctx, hianimeID, episodeID)
	if err != nil {
		return ScrapedStreamMetadata{}, err
	}

	var serverID string
	for _, server := range servers {
		if server.Type == streamType {
			serverID = server.ServerID
			break
		} else if streamType == "sub" && server.Type == "raw" {
			serverID = server.ServerID
			break
		}
	}

	headers := map[string]string{
		"Referer":          c.fetcher.baseURL,
		"X-Requested-With": "XMLHttpRequest",
	}

	var data struct {
		Link string `json:"link"`
	}
	ok, err := c.fetcher.GetAjax(ctx, "/ajax/v2/episode/sources?id="+serverID, headers, &data)
	if err != nil {
		return ScrapedStreamMetadata{}, err
	}
	if !ok {
		return ScrapedStreamMetadata{}, fmt.Errorf("failed to fetch sources")
	}

	parsedURL, err := url.Parse(data.Link)
	if err != nil {
		return ScrapedStreamMetadata{}, err
	}

	parts := strings.Split(parsedURL.Path, "/")
	xrax := parts[len(parts)-1]
	parsedURL.Path = strings.Join(parts[:len(parts)-1], "/")
	parsedURL.RawQuery = ""
	origin := parsedURL.String()

	token, err := c.extractToken(ctx, data.Link)
	if err != nil {
		return ScrapedStreamMetadata{}, fmt.Errorf("failed to extract token: %w", err)
	}

	ajaxURL := fmt.Sprintf("%s/getSources?id=%s&_k=%s", origin, xrax, token)

	ajaxReq, _ := http.NewRequestWithContext(ctx, "GET", ajaxURL, nil)
	ajaxResp, err := c.fetcher.Client.Do(ajaxReq)
	if err != nil {
		return ScrapedStreamMetadata{}, err
	}
	defer ajaxResp.Body.Close()

	var enc ScrapedStreamMetadata
	if err := json.NewDecoder(ajaxResp.Body).Decode(&enc); err != nil {
		return ScrapedStreamMetadata{}, err
	}

	return enc, nil
}

var (
	// 4. window.<key> = "value";
	reWindowString = regexp.MustCompile(`window\.(\w+)\s*=\s*["']([\w-]+)["']`)
	// 5. window.<key> = { ... };
	reWindowObject = regexp.MustCompile(`window\.(\w+)\s*=\s*(\{[\s\S]*?\});`)
	// 5b. extract all string literals inside an object literal
	reQuotedString = regexp.MustCompile(`["']([^"']+)["']`)
)

func (c *HianimeCatalog) extractToken(ctx context.Context, url string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Referer", "https://hianimez.to/")
	resp, err := c.fetcher.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	htmlStr := string(rawBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}

	if meta, exists := doc.Find(`meta[name="_gg_fb"]`).Attr("content"); exists && meta != "" {
		return meta, nil
	}

	if dpi, exists := doc.Find(`[data-dpi]`).Attr("data-dpi"); exists && dpi != "" {
		return dpi, nil
	}

	foundNonce := ""
	doc.Find("script[nonce]").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "empty nonce script") {
			foundNonce, _ = s.Attr("nonce")
			return false // break loop
		}
		return true
	})
	if foundNonce != "" {
		return foundNonce, nil
	}

	if m := reWindowString.FindAllStringSubmatch(htmlStr, -1); len(m) > 0 {
		// take the first match’s value (sub‑match 2)
		return m[0][2], nil
	}

	if m := reWindowObject.FindAllStringSubmatch(htmlStr, -1); len(m) > 0 {
		for _, sub := range m {
			objLiteral := sub[2]
			allStrings := reQuotedString.FindAllStringSubmatch(objLiteral, -1)
			if len(allStrings) == 0 {
				continue
			}
			var sb strings.Builder
			for _, sm := range allStrings {
				sb.WriteString(sm[1])
			}
			if token := sb.String(); len(token) >= 20 {
				return token, nil
			}
		}
	}

	tokenizer := html.NewTokenizer(strings.NewReader(htmlStr))
	const key = "_is_th:"
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return "", nil
		case html.CommentToken:
			data := strings.TrimSpace(string(tokenizer.Text()))
			if strings.HasPrefix(data, key) {
				return strings.TrimPrefix(data, key), nil
			}
		}
	}
}
