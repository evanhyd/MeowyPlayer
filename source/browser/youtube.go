package browser

func NewYouTubeSearcher() Searcher {
	return newClipzagScraper()
}

func NewYouTubeDownloader() Downloader {
	return newMultiDownloader(newCnvmp3Downloader(), newY2APIDownloader(), newY2MateDownloader())
}
