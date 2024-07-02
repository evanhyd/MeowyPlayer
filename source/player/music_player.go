package player

type MusicPlayer interface {
	play()
	prev()
	next()
	setProgress(float64)
	setVolume(float64)
}
