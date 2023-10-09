package cbinding

import (
	"meowyplayer.com/utility/network/scraper"
)

type VideoScraperDataList = dataList[scraper.VideoScraper]

func MakeYoutubeResultDataList() VideoScraperDataList {
	return makeDataList[scraper.VideoScraper]()
}
