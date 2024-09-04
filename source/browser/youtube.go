package browser

func NewYouTubeSearcher() Searcher {
	return newClipzagScraper()
}

func NewYouTubeDownloader() Downloader {
	return newMultiDownloader(newY2APIDownloader(), newY2MateDownloader())
}
