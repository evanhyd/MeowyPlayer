package cbinding

import (
	"golang.org/x/exp/maps"
	"meowyplayer.com/core/resource"
)

type MusicDataList struct {
	dataList[resource.Music]
}

func MakeMusicDataList() MusicDataList {
	return MusicDataList{MakeDataList[resource.Music]()}
}

func (m *MusicDataList) Notify(album resource.Album) {
	m.dataList.Notify(maps.Values(album.MusicList))
}

func (m *MusicDataList) MusicList() []resource.Music {
	return m.data
}
