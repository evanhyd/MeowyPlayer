package scraper

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/resource"
)

var clipzagResultRegex *regexp.Regexp

func init() {
	const clipzagResultPattern = `<a class='title-color' href='(.+)'>\n.+\n.+data-thumb='//(.+)' .+<span class='duration'>(.+)</span></div>\n.+title='(.+)'.+\n.+\n.+\n.+<a class='by-user' href='.+'>(.+)</a><br />(.+)</span>\n.+\n<div class='postdiscription'>(.+)</div>`

	var err error
	clipzagResultRegex, err = regexp.Compile(clipzagResultPattern)
	if err != nil {
		log.Fatal(err)
	}
}

func GetSearchResult(videoTitle string) ([]ClipzagResult, error) {
	//fetch webpage
	const clipzagUrl = `https://clipzag.com/search?`
	queryUrl := clipzagUrl + url.Values{"q": {videoTitle}}.Encode()
	log.Printf("Fetching from: %v\n", queryUrl)
	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//parse result format
	log.Printf("Parsing...")
	content := strings.Builder{}
	if _, err := io.Copy(&content, resp.Body); err != nil {
		return nil, err
	}
	parsed := clipzagResultRegex.FindAllStringSubmatch(content.String(), -1)

	log.Printf("Downloading thumbnails: %v\n", len(parsed))
	results := make([]ClipzagResult, len(parsed))
	completes := make(chan struct{}, len(parsed))
	for i := range parsed {
		go parseSearchResult(parsed[i], &results[len(parsed)-i-1], completes)
	}
	for i := 0; i < len(parsed); i++ {
		<-completes
	}

	log.Printf("Completed")
	return results, nil
}

func parseSearchResult(parsed []string, result *ClipzagResult, completes chan struct{}) {
	staticResource, err := fyne.LoadResourceFromURLString(`https://` + parsed[2])
	if err != nil {
		log.Fatal(err)
	}
	staticImage := canvas.NewImageFromResource(staticResource)
	staticImage.SetMinSize(resource.GetThumbnailIconSize())

	*result = ClipzagResult{
		videoID:     parsed[1][8:],
		thumbnail:   staticImage,
		videoTitle:  html.UnescapeString(parsed[4]),
		stats:       fmt.Sprintf("[%v] %v | %v", parsed[3], parsed[5], parsed[6]),
		description: parsed[7],
	}
	completes <- struct{}{}
}
