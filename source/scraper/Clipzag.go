package scraper

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/PuerkitoBio/goquery"
	"meowyplayer.com/source/resource"
)

func GetSearchResult(videoTitle string) ([]ClipzagResult, error) {
	// Fetch webpage
	const clipzagURL = "https://clipzag.com/search?"
	queryURL := clipzagURL + url.Values{"q": {videoTitle}}.Encode()
	log.Printf("Fetching from: %v", queryURL)
	resp, err := http.Get(queryURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse results
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	posts := doc.Find(".videopost")

	results := make([]ClipzagResult, posts.Length())
	wg := sync.WaitGroup{}
	posts.Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(i int, s *goquery.Selection) {
			defer wg.Done()

			id, _ := s.Find("a.title-color").Attr("href")
			thumb, _ := s.Find("img.videosthumbs-style").Attr("src")
			duration := s.Find(".duration").Text()
			title := s.Find(".title-style").Text()
			user := s.Find("a.by-user").Text()
			stats := s.Find(".viewsanduser span").Contents().Last().Text()
			description := s.Find(".postdiscription").Text()

			imgURL := "https:" + thumb
			staticImage, err := fyne.LoadResourceFromURLString(imgURL)
			if err != nil {
				log.Printf("Error loading thumbnail %v: %v", imgURL, err)
				return
			}

			img := canvas.NewImageFromResource(staticImage)
			img.SetMinSize(resource.GetThumbnailIconSize())

			result := ClipzagResult{
				videoID:     id[8:],
				thumbnail:   img,
				videoTitle:  title,
				stats:       fmt.Sprintf("[%v] %v | %v", duration, user, stats),
				description: description,
			}
			results[posts.Length()-i-1] = result
		}(i, s)
	})

	wg.Wait()

	return results, nil
}
