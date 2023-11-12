package fileformat

import (
	"time"

	"fyne.io/fyne/v2"
)

type VideoResult struct {
	Platform     string
	VideoID      string
	ChannelID    string
	ChannelTitle string
	Title        string
	Stats        string
	Description  string
	Length       time.Duration
	Thumbnail    fyne.Resource
}
