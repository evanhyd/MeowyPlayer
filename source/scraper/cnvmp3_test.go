package scraper

import (
	"testing"
)

func Test_cnvmp3(t *testing.T) {
	testDownload(newCnvmp3Downloader(), &Result{Title: "Renai Circulation", ID: "auQxNYJ07Lc"}, t)
}
