package scraper

import (
	"io"
	"net/http"
	"net/url"
	"regexp"

	"meowyplayer.com/utility/assert"
)

type ClipzagScraper struct {
	regex *regexp.Regexp
}

func (s *ClipzagScraper) NewClipzagScraper() *ClipzagScraper {
	const resultPattern = `<a class="title-color" href="watch\?v=(.+)">\n` + //videoID
		`<div class="video-thumbs">\n` +
		`<img class="videosthumbs-style" data-thumb-m data-thumb="//(.+)" src="//.+"><span class="duration">(.+)</span></div>\n` + //thumbnail, duration
		`<div class="title-style" title="(.+)">.+</div>\n` + //title
		`</a>\n` +
		`<div class="viewsanduser">\n` +
		`<span style="font-weight:bold;"><a class="by-user" href="/channel\?id=(.+)">(.+)</a><br/>(.+)</span>\n` + //channel id, channel title, stats
		`</div>\n` +
		`<div class="postdiscription">(.+)</div>` //description

	regex, err := regexp.Compile(resultPattern)
	assert.NoErr(err)
	return &ClipzagScraper{regex}
}

func (s *ClipzagScraper) getURL(title string) string {
	const clipzagUrl = `https://clipzag.com/search?`
	return clipzagUrl + url.Values{"q": {title}}.Encode()
}

func (s *ClipzagScraper) scrapeContent(content string) ([]ScrapedResult, error) {

	return nil, nil
}

func (s *ClipzagScraper) Search(title string) ([]ScrapedResult, error) {

	//get the response and parse into string
	resp, err := http.Get(s.getURL(title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//scrape
	result, err := s.scrapeContent(string(data))
	if err != nil {
		return nil, err
	}

	return result, nil
}
