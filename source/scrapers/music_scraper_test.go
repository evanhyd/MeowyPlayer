package scrapers

import (
	"context"
	"io"
	"testing"
)

// Helper: test the downloader
func testDownload(t *testing.T, d MusicDownloader, video Result) {
	t.Helper()

	body, err := d.Download(context.Background(), video)
	if err != nil {
		t.Fatalf("failed to start downloading: %v", err)
	}

	if _, err := io.Copy(io.Discard, body); err != nil {
		t.Fatalf("failed to download: %v", err)
	}
	defer body.Close()
}

// Helper: test the searcher
func testSearch(t *testing.T, s MusicSearcher, query string) []Result {
	t.Helper()

	results, err := s.Search(context.Background(), query)
	if err != nil {
		t.Fatalf("Search(%q) failed: %v", query, err)
	}
	if len(results) == 0 {
		t.Fatalf("Search(%q) returned no results", query)
	}
	return results
}

// Table-driven tests for the searcher
func TestPipedSearcher(t *testing.T) {
	searcher := NewPipedSearcher()

	titles := []string{
		"chicken nugget",
		"renai circulation",
		"lofi hip hop",
	}

	for i := range titles {
		t.Run(titles[i], func(t *testing.T) {
			results := testSearch(t, searcher, titles[i])
			t.Logf("Fetched %d results for %q", len(results), titles[i])
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
