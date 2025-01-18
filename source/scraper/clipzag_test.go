package scraper

import (
	"testing"
)

//race condition detection, panic if occurs
//go test -race -run NameOfThatTestFunc .

func TestChickenNugget(t *testing.T) {
	searchQuery(newClipzagScraper(), "chicken nugget", t)
}

func TestMonogatari(t *testing.T) {
	searchQuery(newClipzagScraper(), "renai circulation", t)
}

func searchQuery(scraper *clipzagScraper, title string, t *testing.T) {
	results, err := scraper.Search(title)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	if len(results) == 0 {
		t.Fatalf("failed to fetch results\n")
	}
}
