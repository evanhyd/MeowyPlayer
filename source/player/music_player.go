package player

type MusicPlayer interface {
	Play()
	Prev()
	Next()
	SetProgress(float64)
	SetVolume(float64)
}
