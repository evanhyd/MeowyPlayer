package browser

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

type clipzagScraper struct {
	regex *regexp.Regexp
}

func newClipzagScraper() *clipzagScraper {
	pattern := `<a class='title-color' href='watch\?v=(.+?)'>\s*` + // videoID
		`<div class='video-thumbs'>\s*` +
		`<img class='videosthumbs-style' data-thumb-m(?:='.+?')? data-thumb='//(.+?)' src='//.+?'><span class='duration'>(.+?)</span></div>\s*` + // thumbnail, length
		`<div class='title-style' title='(.+?)'>.+?</div>\s*` + // title
		`</a>\s*` +
		`<div class='viewsanduser'>\s*` +
		`<span style='font-weight:bold;'><a class='by-user' href='/channel\?id=(.+?)'>(.+?)</a><br/>(.+?)</span>\s*` + // channel id, channel title, stats
		`</div>\s*` +
		`<div class='postdiscription'>(.+?)</div>` // description

	return &clipzagScraper{regexp.MustCompile(pattern)}
}

func (s *clipzagScraper) Search(title string) ([]Result, error) {
	page, err := s.fetchPage(title)
	if err != nil {
		return nil, err
	}
	return s.scrapePage(page)
}

func (s *clipzagScraper) fetchPage(title string) (string, error) {
	url := `https://clipzag.com/search?` + url.Values{"q": {title}, "order": {"relevance"}}.Encode()
	resp, err := http.Get(url)
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

func (s *clipzagScraper) scrapePage(content string) ([]Result, error) {
	//parse regex and prepare output buffers
	matches := s.regex.FindAllStringSubmatch(content, -1)
	results := make([]Result, len(matches))
	errors := make(chan error, len(matches))

	//parse into the results
	wg := sync.WaitGroup{}
	wg.Add(len(matches))

	go func() {
		for i := range matches {
			go func(match []string, result *Result) {
				defer wg.Done()
				s.parseMatch(match, result, errors)
			}(matches[i], &results[i])
		}
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (s *clipzagScraper) parseMatch(match []string, result *Result, errors chan<- error) {
	//download thumbnail
	thumbnail, err := fyne.LoadResourceFromURLString(`https://` + match[2])
	if err != nil {
		errors <- err
	}

	//calculate video length
	hourMinSec := strings.Split(match[3], ":")
	totalSecond := int64(0)
	for _, time := range hourMinSec {
		t, err := strconv.ParseInt(time, 10, 64)
		if err != nil {
			errors <- err
			return
		}
		totalSecond = totalSecond*60 + t
	}

	*result = Result{
		Platform:     "YouTube",
		ID:           match[1],
		Thumbnail:    thumbnail,
		Length:       time.Duration(totalSecond * int64(time.Second)),
		Title:        html.UnescapeString(match[4]),
		ChannelID:    match[5],
		ChannelTitle: html.UnescapeString(match[6]),
		Stats:        html.UnescapeString(match[7]),
		Description:  html.UnescapeString(match[8]),
	}
}
