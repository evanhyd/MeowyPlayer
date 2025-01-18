package scraper

import (
	"testing"
)

func testDownload(downloader Downloader, video *Result, t *testing.T) {
	body, err := downloader.Download(video)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer body.Close()

	// file, err := os.Create("test.mp3")
	// if err != nil {
	// 	t.Fatalf("%v\n", err)
	// }
	// defer file.Close()
	// if _, err := io.Copy(file, body); err != nil {
	// 	t.Fatalf("%v\n", err)
	// }
}
