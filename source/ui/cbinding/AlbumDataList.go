package cbinding

import "meowyplayer.com/source/player"

type AlbumDataList struct {
	dataList[player.Album]
}

func MakeAlbumDataList() AlbumDataList {
	return AlbumDataList{makeDataList[player.Album]()}
}

func (a *AlbumDataList) Notify(config *player.Config) {
	a.dataList.Notify(config.Albums)
}
