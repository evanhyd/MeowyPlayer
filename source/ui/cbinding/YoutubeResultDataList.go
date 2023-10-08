package cbinding

import (
	"meowyplayer.com/utility/network/scraper"
)

type YoutubeResultDataList = dataList[scraper.YoutubeResult]

func MakeYoutubeResultDataList() YoutubeResultDataList {
	return makeDataList[scraper.YoutubeResult]()
}
