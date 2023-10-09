package scraper

import (
	"time"

	"fyne.io/fyne/v2"
)

type VideoResult struct {
	VideoID      string
	ChannelID    string
	ChannelTitle string
	Title        string
	Stats        string
	Description  string
	Length       time.Duration
	Thumbnail    fyne.Resource
}

type VideoScraper interface {
	Search(title string) ([]VideoResult, error)
}
