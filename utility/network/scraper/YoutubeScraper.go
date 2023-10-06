package scraper

import (
	"time"

	"fyne.io/fyne/v2"
)

type YoutubeResult struct {
	VideoID      string
	Thumbnail    fyne.Resource
	Length       time.Duration
	Title        string
	ChannelID    string
	ChannelTitle string
	Stats        string
	Description  string
}

type YoutubeScraper interface {
	Search(title string) ([]YoutubeResult, error)
}
