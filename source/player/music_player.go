package player

import (
	"meowyplayer/storage"
)

type QueueMode = int64

const (
	SequentialQueueMode QueueMode = iota
	RandomQueueMode
)

type MusicPlayer interface {
	Resume()
	Suspend()
	Previous()
	Next()
	SetQueueMode(mode QueueMode)
	SetRepeat(isRepeating bool)
	SetProgress(percent float64)
	SetVolume(percent float64)

	// Event Handler
	OnSelectPlaylist(musicList []storage.Music, index int)
}
