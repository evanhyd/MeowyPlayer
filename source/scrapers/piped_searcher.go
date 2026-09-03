package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"meowyplayer/storages"
)

var _ MusicSearcher = (*pipedSearcher)(nil)

type PipedVideoItem struct {
	URL              string `json:"url"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	UploaderName     string `json:"uploaderName"`
	UploaderURL      string `json:"uploaderUrl"`
	Duration         int64  `json:"duration"`
	Views            int64  `json:"views"`
	ShortDescription string `json:"shortDescription"`
	IsShort          bool   `json:"isShort"`
}

type PipedSearchResponse struct {
	Items []PipedVideoItem `json:"items"`
}

type pipedSearcher struct {
	apiURL string
}

func NewPipedSearcher() *pipedSearcher {
	return &pipedSearcher{apiURL: "https://api.piped.private.coffee"}
}

func (s *pipedSearcher) Search(ctx context.Context, title string) ([]Result, error) {
	endpoint := fmt.Sprintf("%s/search?q=%s&filter=videos", s.apiURL, url.QueryEscape(title))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Piped API] error response: %v", resp.Status)
	}

	var searchResp PipedSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("[Piped API] failed to decode JSON (%v): %v", err, body)
	}

	var videoItems []PipedVideoItem
	for _, item := range searchResp.Items {
		if item.Type == "stream" && !item.IsShort && item.Duration > 0 {
			videoItems = append(videoItems, item)
		}
	}

	results := make([]Result, len(videoItems))
	errors := make([]error, len(videoItems))

	var wg sync.WaitGroup
	wg.Add(len(videoItems))

	for i, item := range videoItems {
		go func(index int, vid PipedVideoItem) {
			defer wg.Done()
			results[index], errors[index] = s.parseItem(ctx, vid)
		}(i, item)
	}
	wg.Wait()

	var validResults []Result
	for i, err := range errors {
		if err != nil {
			if ctx.Err() == nil {
				vidID := strings.TrimPrefix(videoItems[i].URL, "/watch?v=")
				fmt.Printf("Skipping video %s due to error: %v\n", vidID, err)
			}
			continue
		}
		validResults = append(validResults, results[i])
	}

	return validResults, nil
}

func (s *pipedSearcher) parseItem(ctx context.Context, item PipedVideoItem) (Result, error) {
	videoID := strings.TrimPrefix(item.URL, "/watch?v=")
	channelID := strings.TrimPrefix(item.UploaderURL, "/channel/")
	thumbURL := fmt.Sprintf("https://i.ytimg.com/vi/%s/mqdefault.jpg", videoID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, thumbURL, nil)
	if err != nil {
		return Result{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("failed to fetch thumbnail: %s", resp.Status)
	}

	thumbnailData, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Platform:     storages.YouTubeSource,
		ID:           videoID,
		Thumbnail:    thumbnailData,
		Length:       time.Duration(item.Duration * int64(time.Second)),
		Title:        item.Title,
		ChannelID:    channelID,
		ChannelTitle: item.UploaderName,
		Stats:        fmt.Sprintf("%d views", item.Views),
		Description:  item.ShortDescription,
	}, nil
}
