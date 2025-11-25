package scrapers

import (
	"io"
	"meowyplayer/storages"
	"time"
)

type Result struct {
	Platform     storages.MusicSource
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

type MusicScraper interface {
	MusicSearcher
	MusicDownloader
}

func NewYouTubeScraper() MusicScraper {
	type youtubeScraper struct {
		MusicSearcher
		MusicDownloader
	}
	return youtubeScraper{newClipzagSearcher(), newCnvmp3Downloader()}
}
