package browser

import (
	"testing"
)

// race condition detection, panic if occurs
// go test -race -run NameOfThatTestFunc .
func Test_y2Mate(t *testing.T) {
	DownloadQuery(newY2MateDownloader(), &Result{Title: "Renai Circulation", ID: "auQxNYJ07Lc"}, t)
}
