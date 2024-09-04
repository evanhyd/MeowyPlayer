package browser

import "testing"

func DownloadQuery(downloader Downloader, video *Result, t *testing.T) {
	_, err := downloader.Download(video)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
