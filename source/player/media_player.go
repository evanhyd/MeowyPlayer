package player

type MediaPlayer interface {
	Play()
	Prev()
	Next()
	SetProgress(float64)
	SetVolume(float64)
}
