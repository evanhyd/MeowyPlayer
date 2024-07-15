package browser

var _ Browser = &youTubeBrowser{}

type youTubeBrowser struct {
	*clipzagScraper
	*y2MateDownloader
}

func NewYouTubeBrowser() *youTubeBrowser {
	return &youTubeBrowser{newClipzagScraper(), newY2MateDownloader()}
}
