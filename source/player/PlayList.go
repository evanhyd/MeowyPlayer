package player

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/utility"
)

type PlayList struct {
	album      Album
	musicIndex int
}

func NewPlayList(album *Album, music *Music) *PlayList {
	index := slices.Index(album.MusicList, *music)
	utility.Assert(func() bool { return index != -1 })
	return &PlayList{*album, index}
}

func (p *PlayList) Album() *Album {
	return &p.album
}

func (p *PlayList) Music() *Music {
	return &p.album.MusicList[p.musicIndex]
}

func (p *PlayList) NextMusic() {
	p.musicIndex = (p.musicIndex + 1) % len(p.album.MusicList)
}

func (p *PlayList) PrevMusic() {
	p.musicIndex = (p.musicIndex - 1 + len(p.album.MusicList)) % len(p.album.MusicList)
}
