package player

import (
	"slices"

	"meowyplayer.com/core/resource"
)

type PlayList struct {
	musicList []resource.Music
	index     int
}

func MakePlayList(musicList []resource.Music, music *resource.Music) PlayList {
	return PlayList{musicList, slices.Index(musicList, *music)}
}

func (p *PlayList) Music() *resource.Music {
	return &p.musicList[p.index]
}

func (p *PlayList) MusicCount() int {
	return len(p.musicList)
}

func (p *PlayList) Index() int {
	return p.index
}

func (p *PlayList) SetIndex(musicIndex int) {
	p.index = musicIndex
}
