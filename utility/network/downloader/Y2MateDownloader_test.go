package downloader_test

import (
	"testing"

	"meowyplayer.com/utility/network/downloader"
	"meowyplayer.com/utility/network/fileformat"
)

//race condition detection, panic if occurs
//go test -race -run NameOfThatTestFunc .

func TestRenaiCirculation(t *testing.T) {
	DownloadQuery(downloader.NewY2MateDownloader(), &fileformat.VideoResult{Title: "Renai Circulation", VideoID: "auQxNYJ07Lc"}, t)
}

func TestMousoExpress(t *testing.T) {
	DownloadQuery(downloader.NewY2MateDownloader(), &fileformat.VideoResult{Title: "Mouso Express", VideoID: "y2XArpEcygc"}, t)
}

func TestIntoTheLight(t *testing.T) {
	DownloadQuery(downloader.NewY2MateDownloader(), &fileformat.VideoResult{Title: "Into The Light", VideoID: "uYO7zbc-wJ0"}, t)
}

func DownloadQuery(downloader downloader.MusicDownloader, video *fileformat.VideoResult, t *testing.T) {
	_, err := downloader.Download(video)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
