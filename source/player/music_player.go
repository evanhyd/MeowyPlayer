package player

import "meowyplayer/storage"

type QueueMode = int64

const (
	Order QueueMode = iota
	Random
)

type MusicPlayer interface {
	Resume()
	Suspend()
	Previous()
	Next()
	SetQueueMode(QueueMode)
	SetRepeat(bool)
	SetProgress(float64)
	SetVolume(float64)
	SetPlaylist([]storage.Music)
}
