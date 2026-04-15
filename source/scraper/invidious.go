package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// invSearchItem strictly matches the Invidious JSON payload you provided
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

func newInvidiousSearcher() *invidiousSearcher {
	return &invidiousSearcher{
		apiBaseURL: "https://inv.thepixora.com",
	}
}

func (s *invidiousSearcher) Search(title string) ([]Result, error) {
	data, err := s.fetchSearchPage(title)
	if err != nil {
		return nil, err
	}
	return s.scrapeSearchPage(data)
}

func (s *invidiousSearcher) fetchSearchPage(title string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/api/v1/search?q=%s&type=video", s.apiBaseURL, url.QueryEscape(title))

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Invidious API] error response: %v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	return data, err
}

func (s *invidiousSearcher) scrapeSearchPage(content []byte) ([]Result, error) {
	var items []invSearchItem
	if err := json.Unmarshal(content, &items); err != nil {
		return nil, fmt.Errorf("[Invidious API] failed to decode JSON: %v", err)
	}

	// Filter out any accidental channels or playlists
	var videoItems []invSearchItem
	for _, item := range items {
		if item.Type == "video" {
			videoItems = append(videoItems, item)
		}
	}

	// Prepare output buffers
	results := make([]Result, len(videoItems))
	errors := make(chan error, len(videoItems))

	// Parse into the results concurrently
	wg := sync.WaitGroup{}
	wg.Add(len(videoItems))
	go func() {
		for i := range videoItems {
			go func(item invSearchItem, result *Result) {
				defer wg.Done()
				s.parseMatchResult(item, result, errors)
			}(videoItems[i], &results[i])
		}
		wg.Wait()
		close(errors)
	}()

	// Fail fast if any errors occurred during concurrent processing
	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (s *invidiousSearcher) parseMatchResult(item invSearchItem, result *Result, errors chan<- error) {
	// Download thumbnail using fyne
	thumbURL := fmt.Sprintf("https://i.ytimg.com/vi/%s/mqdefault.jpg", item.VideoID)
	thumbnail, err := fyne.LoadResourceFromURLString(thumbURL)
	if err != nil {
		errors <- err
		return
	}

	// Map to result struct directly via pointer
	*result = Result{
		Platform:     "YouTube",
		ID:           item.VideoID,
		Thumbnail:    thumbnail,
		Length:       time.Duration(item.LengthSeconds * int64(time.Second)),
		Title:        item.Title,
		ChannelID:    item.AuthorID,
		ChannelTitle: item.Author,
		Stats:        item.ViewCountText,
		Description:  item.Description,
	}
}
