package player

import (
	"meowyplayer.com/source/pattern"
)

var state State

func init() {
	state = State{}
}

func GetState() *State {
	return &state
}

type State struct {
	album           Album
	musics          []Music
	onUpdateAlbums  pattern.OneArgObservable[[]Album]
	onUpdateMusics  pattern.OneArgObservable[[]Music]
	onUpdateSeeker  pattern.ThreeArgObservable[Album, []Music, Music]
	onFocusAlbumTab pattern.ZeroArgObservable
	onFocusMusicTab pattern.ZeroArgObservable
}

func (state *State) Album() Album {
	return state.album
}

func (state *State) OnUpdateAlbums() pattern.OneArgObservabler[[]Album] {
	return &state.onUpdateAlbums
}

func (state *State) OnUpdateMusics() pattern.OneArgObservabler[[]Music] {
	return &state.onUpdateMusics
}

func (state *State) OnUpdateSeeker() pattern.ThreeArgObservabler[Album, []Music, Music] {
	return &state.onUpdateSeeker
}

func (state *State) OnFocusAlbumTab() pattern.ZeroArgObservabler {
	return &state.onFocusAlbumTab
}

func (state *State) OnFocusMusicTab() pattern.ZeroArgObservabler {
	return &state.onFocusMusicTab
}
