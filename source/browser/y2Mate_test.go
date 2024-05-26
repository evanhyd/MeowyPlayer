package browser

import (
	"testing"
)

//race condition detection, panic if occurs
//go test -race -run NameOfThatTestFunc .

func TestRenaiCirculation(t *testing.T) {
	DownloadQuery(newY2MateDownloader(), &Result{Title: "Renai Circulation", VideoID: "auQxNYJ07Lc"}, t)
}

func TestMousoExpress(t *testing.T) {
	DownloadQuery(newY2MateDownloader(), &Result{Title: "Mouso Express", VideoID: "y2XArpEcygc"}, t)
}

func TestIntoTheLight(t *testing.T) {
	DownloadQuery(newY2MateDownloader(), &Result{Title: "Into The Light", VideoID: "uYO7zbc-wJ0"}, t)
}

func TestLoveStory(t *testing.T) {
	DownloadQuery(newY2MateDownloader(), &Result{Title: "Love Story", VideoID: "8xg3vE8Ie_E"}, t)
}

func DownloadQuery(downloader *y2MateDownloader, video *Result, t *testing.T) {
	_, err := downloader.Download(video)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
