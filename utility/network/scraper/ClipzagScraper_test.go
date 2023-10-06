package scraper_test

import (
	"testing"

	"meowyplayer.com/utility/network/scraper"
)

func TestChickenNugget(t *testing.T) {
	var scraper scraper.Scraper = &scraper.ClipzagScraper{}

	_, err := scraper.Search("chicken nugget")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
