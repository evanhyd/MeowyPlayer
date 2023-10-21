package downloader

import "meowyplayer.com/utility/network/fileformat"

type MusicDownloader interface {
	Download(video *fileformat.VideoResult) ([]byte, error)
}
