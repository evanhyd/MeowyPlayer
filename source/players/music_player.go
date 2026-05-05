package players

import "meowyplayer/storages"

type QueueMode = int64

const (
	SequentialQueueMode QueueMode = iota
	RandomQueueMode
)

type MusicPlayer interface {
	IsPlaying() bool
	GetProgress() float64
	GetMusic() storages.Music

	Resume()
	Pause()
	Previous()
	Next()

	SetPlaylist(playlist storages.Playlist, selectedMusic storages.Music)
	SetQueueMode(mode QueueMode)
	SetRepeat(isRepeating bool)
	SetProgress(percent float64)
	SetVolume(percent float64)
}
