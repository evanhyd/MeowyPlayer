package scrapers

import (
	"testing"
)

// Helper: test the downloader
func testDownload(t *testing.T, d MusicDownloader, video Result) {
	t.Helper()

	body, err := d.Download(video)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	defer body.Close()
}

// Helper: test the searcher
func testSearch(t *testing.T, s MusicSearcher, query string) []Result {
	t.Helper()

	results, err := s.Search(query)
	if err != nil {
		t.Fatalf("Search(%q) failed: %v", query, err)
	}
	if len(results) == 0 {
		t.Fatalf("Search(%q) returned no results", query)
	}
	return results
}

// Table-driven tests for the searcher
func TestClipzagSearcher(t *testing.T) {
	searcher := NewClipzagSearche()

	tests := []string{
		"chicken nugget",
		"renai circulation",
		"lofi hip hop",
	}

	for _, q := range tests {
		t.Run(q, func(t *testing.T) {
			results := testSearch(t, searcher, q)
			t.Logf("Fetched %d results for %q", len(results), q)
		})
	}
}

// Test the cnvmp3 downloader directly
func TestCnvmp3Downloader(t *testing.T) {
	video := Result{
		ID:    "auQxNYJ07Lc",
		Title: "Renai Circulation",
	}
	testDownload(t, NewCnvmp3Downloader(), video)
}
