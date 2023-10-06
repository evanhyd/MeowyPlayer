package scraper

import "fyne.io/fyne/v2"

type ScrapedResult struct {
	Thumbnail   fyne.Resource
	VideoID     string
	Title       string
	Description string
}

type Scraper interface {
	Search(title string) ([]ScrapedResult, error)
}
