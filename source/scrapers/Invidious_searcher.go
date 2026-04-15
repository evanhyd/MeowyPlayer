package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"meowyplayer/storages"
)

type invSearchItem struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	VideoID       string `json:"videoId"`
	Author        string `json:"author"`
	AuthorID      string `json:"authorId"`
	LengthSeconds int64  `json:"lengthSeconds"`
	ViewCountText string `json:"viewCountText"`
	Description   string `json:"description"`
}

type invidiousSearcher struct {
	apiBaseURL string
}

func NewInvidiousSearcher() *invidiousSearcher {
	return &invidiousSearcher{
		apiBaseURL: "https://inv.thepixora.com",
	}
}

func (s *invidiousSearcher) Search(ctx context.Context, title string) ([]Result, error) {
	endpoint := fmt.Sprintf("%s/api/v1/search?q=%s&type=video", s.apiBaseURL, url.QueryEscape(title))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Invidious API] error response: %v", resp.Status)
	}

	var items []invSearchItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("[Invidious API] failed to decode JSON: %v", err)
	}

	// Filter out any accidental channels or playlists that might slip into the search
	var videoItems []invSearchItem
	for _, item := range items {
		if item.Type == "video" {
			videoItems = append(videoItems, item)
		}
	}

	results := make([]Result, len(videoItems))
	errors := make([]error, len(videoItems))

	// Parse the results and fetch thumbnails concurrently
	var wg sync.WaitGroup
	wg.Add(len(videoItems))

	for i, item := range videoItems {
		go func(index int, vid invSearchItem) {
			defer wg.Done()
			results[index], errors[index] = s.parseItem(ctx, vid)
		}(i, item)
	}
	wg.Wait()

	// Filter out any results that failed to download their thumbnail
	var validResults []Result
	for i, err := range errors {
		if err != nil {
			// Only log errors if the user didn't intentionally cancel the search
			if ctx.Err() == nil {
				fmt.Printf("Skipping video %s due to error: %v\n", videoItems[i].VideoID, err)
			}
			continue
		}
		validResults = append(validResults, results[i])
	}

	return validResults, nil
}

func (s *invidiousSearcher) parseItem(ctx context.Context, item invSearchItem) (Result, error) {
	// Hitting YouTube's image server directly is significantly faster.
	thumbURL := fmt.Sprintf("https://i.ytimg.com/vi/%s/mqdefault.jpg", item.VideoID)

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
		ID:           item.VideoID,
		Thumbnail:    thumbnailData,
		Length:       time.Duration(item.LengthSeconds * int64(time.Second)),
		Title:        item.Title,
		ChannelID:    item.AuthorID,
		ChannelTitle: item.Author,
		Stats:        item.ViewCountText,
		Description:  item.Description,
	}, nil
}
