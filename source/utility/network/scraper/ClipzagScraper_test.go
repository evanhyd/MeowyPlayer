package scraper_test

import (
	"testing"

	"meowyplayer.com/utility/network/scraper"
)

//race condition detection, panic if occurs
//go test -race -run NameOfThatTestFunc .

func TestChickenNugget(t *testing.T) {
	SearchQuery(scraper.NewClipzagScraper(), "chicken nugget", t)
}

func TestMonogatari(t *testing.T) {
	SearchQuery(scraper.NewClipzagScraper(), "renai circulation", t)
}

func SearchQuery(scraper scraper.VideoScraper, title string, t *testing.T) {
	results, err := scraper.Search(title)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	for _, result := range results {
		t.Logf("\n\n%+v\n\n", result)
	}
}
