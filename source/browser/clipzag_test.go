package browser

import (
	"testing"
)

//race condition detection, panic if occurs
//go test -race -run NameOfThatTestFunc .

func TestChickenNugget(t *testing.T) {
	SearchQuery(newClipzagScraper(), "chicken nugget", t)
}

func TestMonogatari(t *testing.T) {
	SearchQuery(newClipzagScraper(), "renai circulation", t)
}

func SearchQuery(scraper *clipzagScraper, title string, t *testing.T) {
	results, err := scraper.Search(title)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	if len(results) == 0 {
		t.Fatalf("failed to fetch results\n")
	}
}
