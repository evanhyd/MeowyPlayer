package scraper

import "meowyplayer.com/utility/network/fileformat"

type VideoScraper interface {
	Search(title string) ([]fileformat.VideoResult, error)
}
