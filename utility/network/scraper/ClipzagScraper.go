package scraper

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"meowyplayer.com/utility/assert"
)

type ClipzagScraper struct {
	regex *regexp.Regexp
}

func NewClipzagScraper() *ClipzagScraper {
	const resultPattern = `<a class="title-color" href="watch\?v=(.+)">\n` + //videoID
		`<div class="video-thumbs">\n` +
		`<img class="videosthumbs-style" data-thumb-m data-thumb="//(.+)" src="//.+"><span class="duration">(.+)</span></div>\n` + //thumbnail, length
		`<div class="title-style" title="(.+)">.+</div>\n` + //title
		`</a>\n` +
		`<div class="viewsanduser">\n` +
		`<span style="font-weight:bold;"><a class="by-user" href="/channel\?id=(.+)">(.+)</a><br/>(.+)</span>\n` + //channel id, channel title, stats
		`</div>\n` +
		`<div class="postdiscription">(.+)</div>` //description

	regex, err := regexp.Compile(resultPattern)
	assert.NoErr(err, "failed to compile Clipzag scraper regex")
	return &ClipzagScraper{regex}
}

func (s *ClipzagScraper) Search(title string) ([]YoutubeResult, error) {
	content, err := s.getResponse(title)
	if err != nil {
		return nil, err
	}
	return s.scrapeContent(content), nil
}

func (s *ClipzagScraper) getResponse(title string) (string, error) {
	url := `https://clipzag.com/search?` + url.Values{"q": {title}, "order": {"relevance"}}.Encode()
	log.Printf("scraping from %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	builder := &strings.Builder{}
	if _, err := io.Copy(builder, resp.Body); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (s *ClipzagScraper) scrapeContent(content string) []YoutubeResult {
	//parse regex and prepare output buffers
	matches := s.regex.FindAllStringSubmatch(content, -1)
	results := make([]YoutubeResult, len(matches))
	log.Printf("scraping %v results...\n", len(matches))

	//parse into the results
	wg := sync.WaitGroup{}
	wg.Add(len(matches))
	for i, match := range matches {
		i := i
		match := match
		go func() {
			defer wg.Done()
			s.parseMatch(match, &results[i])
		}()
	}
	wg.Wait()

	log.Println("scraping completed")
	return results

	//This code doesn't work well, since the output of the channel is not in order
	//This can cause unwanted result appear on the top of the search results
	/*
		//parse matches into results concurrently
		matches := s.regex.FindAllStringSubmatch(content, -1)
		log.Printf("scraping %v results...\n", len(matches))
		resultChan := make(chan YoutubeResult, len(matches))
		for _, match := range matches {
			go s.parseMatch(match, resultChan)
		}

		//collect results
		results := make([]YoutubeResult, 0, len(matches))
		for range matches {
			results = append(results, <-resultChan)
		}

		return results
	*/
}

func (s *ClipzagScraper) parseMatch(match []string, dst *YoutubeResult) {
	thumbnail, err := fyne.LoadResourceFromURLString(`https://` + match[2])
	assert.NoErr(err, "failed to download the thumbnail")

	times := strings.Split(match[3], ":")
	seconds := 0
	for _, time := range times {
		t, err := strconv.Atoi(time)
		assert.NoErr(err, "invalid time conversion")
		seconds = seconds*60 + t
	}

	*dst = YoutubeResult{
		VideoID:      match[1],
		Thumbnail:    thumbnail,
		Length:       time.Duration(seconds * int(time.Second)),
		Title:        html.UnescapeString(match[4]),
		ChannelID:    match[5],
		ChannelTitle: html.UnescapeString(match[6]),
		Stats:        html.UnescapeString(match[7]),
		Description:  html.UnescapeString(match[8]),
	}
}
