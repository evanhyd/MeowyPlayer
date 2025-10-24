package scraper

import (
	"io"
	"meowyplayer/storage"
	"time"
)

type Result struct {
	Platform     storage.MusicSource
	ID           string
	ChannelID    string
	ChannelTitle string
	Title        string
	Stats        string
	Description  string
	Length       time.Duration
	Thumbnail    []byte
}

type MusicSearcher interface {
	Search(string) ([]Result, error)
}

type MusicDownloader interface {
	Download(Result) (io.ReadCloser, error)
}

type MusicScraper struct {
	MusicSearcher
	MusicDownloader
}

func NewYouTubeScraper() MusicScraper {
	return MusicScraper{newClipzagSearcher(), newCnvmp3Downloader()}
}
