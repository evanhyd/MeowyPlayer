package cbinding

import "meowyplayer.com/source/resource"

type MusicDataList struct {
	dataList[resource.Music]
	album resource.Album
}

func MakeMusicDataList() MusicDataList {
	return MusicDataList{makeDataList[resource.Music](), resource.Album{}}
}

func (m *MusicDataList) Notify(album *resource.Album) {
	m.album = *album
	m.dataList.Notify(album.MusicList)
}

func (m *MusicDataList) GetAlbum() *resource.Album {
	return &m.album
}
