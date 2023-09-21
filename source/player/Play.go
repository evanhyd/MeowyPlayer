package player

import (
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/utility"
)

type Play struct {
	album      Album
	musicIndex int
}

func NewPlay(album *Album, music *Music) *Play {
	index := slices.Index(album.MusicList, *music)
	utility.Assert(func() bool { return index != -1 })
	return &Play{*album, index}
}

func (p *Play) Album() *Album {
	return &p.album
}

func (p *Play) Music() *Music {
	return &p.album.MusicList[p.musicIndex]
}

func (p *Play) NextMusic() {
	p.musicIndex = (p.musicIndex + 1) % len(p.album.MusicList)
}

func (p *Play) PrevMusic() {
	p.musicIndex = (p.musicIndex - 1 + len(p.album.MusicList)) % len(p.album.MusicList)
}
