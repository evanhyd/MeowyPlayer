package scraper_test

import (
	"testing"

	"meowyplayer.com/utility/scraper"
)

func TestSearch(t *testing.T) {
	title := "monogatari"
	err := scraper.Search(title)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
