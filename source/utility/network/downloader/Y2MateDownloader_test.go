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

func TestLoveStory(t *testing.T) {
	DownloadQuery(downloader.NewY2MateDownloader(), &fileformat.VideoResult{Title: "Love Story", VideoID: "8xg3vE8Ie_E"}, t)
}

func DownloadQuery(downloader downloader.MusicDownloader, video *fileformat.VideoResult, t *testing.T) {
	data, err := downloader.Download(video)
	if len(data) == 0 || err != nil {
		t.Fatalf("%v\n", err)
	}
}
