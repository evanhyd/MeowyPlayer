package scraper

func NewYouTubeSearcher() Searcher {
	return newClipzagScraper()
}

func NewYouTubeDownloader() Downloader {
	return newMultiDownloader(newCnvmp3Downloader())
}
