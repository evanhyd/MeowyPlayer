package browser

import (
	"testing"
)

//race condition detection, panic if occurs
//go test -race -run NameOfThatTestFunc .

func Test_y2API(t *testing.T) {
	DownloadQuery(newY2APIDownloader(), &Result{Title: "Renai Circulation", ID: "auQxNYJ07Lc"}, t)
}
