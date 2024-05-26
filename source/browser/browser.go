package browser

import (
	"io"
	"time"

	"fyne.io/fyne/v2"
)

type Result struct {
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

type Browser interface {
	Search(string) ([]Result, error)
	Download(*Result) (io.ReadCloser, error)
}
