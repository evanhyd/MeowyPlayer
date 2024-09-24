package browser

import (
	"testing"
)

func Test_y2API(t *testing.T) {
	DownloadQuery(newY2APIDownloader(), &Result{Title: "Renai Circulation", ID: "auQxNYJ07Lc"}, t)
}
