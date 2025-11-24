package scrapers

import (
	"fmt"
	"html"
	"io"
	"meowyplayer/storages"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type clipzagSearcher struct {
	matchResultRegex *regexp.Regexp
}

func newClipzagSearcher() *clipzagSearcher {
	matchResultPattern := `<a class='title-color' href='watch\?v=(.+?)'>\s*` + // videoID
		`<div class='video-thumbs'>\s*` +
		`<img class='videosthumbs-style' data-thumb-m(?:='.*?')? data-thumb='//(.+?)' src='//.+?'><span class='duration'>(.+?)</span></div>\s*` + // thumbnail, length
		`<div class='title-style' title='(.+?)'>.+?</div>\s*` + // title
		`</a>\s*` +
		`<div class='viewsanduser'>\s*` +
		`<span style='font-weight:bold;'><a class='by-user' href='/channel\?id=(.+?)'>(.+?)</a><br/>(.+?)</span>\s*` + // channel id, channel title, stats
		`</div>\s*` +
		`<div class='postdiscription'>(.+?)</div>` // description

	return &clipzagSearcher{regexp.MustCompile(matchResultPattern)}
}

func (s *clipzagSearcher) Search(title string) ([]Result, error) {
	page, err := s.fetchSearchPage(title)
	if err != nil {
		return nil, err
	}
	return s.scrapeSearchPage(page)
}

func (s *clipzagSearcher) fetchSearchPage(title string) (string, error) {
	endpoint := `https://clipzag.com/search?` + url.Values{"q": {title}, "order": {"relevance"}}.Encode()
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[Clipzag] error response: %v", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	return string(data), err
}

func (s *clipzagSearcher) scrapeSearchPage(content string) ([]Result, error) {
	// Parse regex and prepare output buffers.
	matches := s.matchResultRegex.FindAllStringSubmatch(content, -1)
	results := make([]Result, len(matches))
	errors := make([]error, len(matches))

	// Parse the result concurrently..
	wg := sync.WaitGroup{}
	wg.Add(len(matches))
	for i := range matches {
		go func() {
			defer wg.Done()
			results[i], errors[i] = s.parseMatchResult(matches[i])
		}()
	}
	wg.Wait()

	// Any error is error.
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (s *clipzagSearcher) parseMatchResult(match []string) (Result, error) {
	// Download the thumbnail.
	resp, err := http.Get(`https://` + match[2])
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	thumbnail, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	// Calculate the video length.
	hourMinSec := strings.Split(match[3], ":")
	totalSecond := int64(0)
	for _, time := range hourMinSec {
		t, err := strconv.ParseInt(time, 10, 64)
		if err != nil {
			return Result{}, err
		}
		totalSecond = totalSecond*60 + t
	}

	return Result{
		Platform:     storages.YouTubeSource,
		ID:           match[1],
		Thumbnail:    thumbnail,
		Length:       time.Duration(totalSecond * int64(time.Second)),
		Title:        html.UnescapeString(match[4]),
		ChannelID:    match[5],
		ChannelTitle: html.UnescapeString(match[6]),
		Stats:        html.UnescapeString(match[7]),
		Description:  html.UnescapeString(match[8]),
	}, nil
}
