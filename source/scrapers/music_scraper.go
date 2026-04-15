package scrapers

import (
	"context"
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
	Search(context.Context, string) ([]Result, error)
}

type MusicDownloader interface {
	Download(context.Context, Result) (io.ReadCloser, error)
}
