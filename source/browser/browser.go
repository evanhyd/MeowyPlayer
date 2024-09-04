package browser

import (
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
)

type Result struct {
	Platform     string
	ID           string
	ChannelID    string
	ChannelTitle string
	Title        string
	Stats        string
	Description  string
	Length       time.Duration
	Thumbnail    fyne.Resource
}
type Searcher interface {
	Search(string) ([]Result, error)
}

type Downloader interface {
	Download(*Result) (io.ReadCloser, error)
}

type MultiDownloader struct {
	downloaders []Downloader
}

func newMultiDownloader(downloaders ...Downloader) *MultiDownloader {
	return &MultiDownloader{downloaders: downloaders}
}

func (d *MultiDownloader) Download(result *Result) (io.ReadCloser, error) {
	var errs []error
	for _, downloader := range d.downloaders {
		content, err := downloader.Download(result)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		return content, nil
	}
	return nil, fmt.Errorf("all downloader failed: %v", errs)
}
