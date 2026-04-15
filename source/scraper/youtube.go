package scraper

func NewYouTubeSearcher() Searcher {
	return newInvidiousSearcher()
}

func NewYouTubeDownloader() Downloader {
	return newMultiDownloader(newCnvmp3Downloader())
}
