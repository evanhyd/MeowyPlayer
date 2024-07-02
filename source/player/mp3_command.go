package player

import "io"

type MP3Command interface {
	execute(*MP3Player)
}

type mp3Play struct{}
type mp3Prev struct{}
type mp3Next struct{}
type mp3Progress struct{ percent float64 }
type m3pVolume struct{ percent float64 }

func (mp3Play) execute(mp3Player *MP3Player) {
	if mp3Player.player.IsPlaying() {
		mp3Player.player.Pause()
	} else {
		mp3Player.player.Play()
	}
}

func (mp3Prev) execute(mp3Player *MP3Player) {
	mp3Player.historyIndex = max(0, mp3Player.historyIndex-1)
	music := mp3Player.history[mp3Player.historyIndex]
	mp3Player.loadMusic(music)
}

func (mp3Next) execute(mp3Player *MP3Player) {
	//reading the history queue
	mp3Player.historyIndex = min(len(mp3Player.history), mp3Player.historyIndex+1)
	if mp3Player.historyIndex < len(mp3Player.history) {
		music := mp3Player.history[mp3Player.historyIndex]
		mp3Player.loadMusic(music)
		return
	}

	switch mp3Player.mode {
	case kRandomMode:
		mp3Player.playIndex = mp3Player.queue[mp3Player.queueIndex]
		mp3Player.queueIndex = (mp3Player.queueIndex + 1) % len(mp3Player.queue)

	case kOrderedMode:
		mp3Player.playIndex = (mp3Player.playIndex + 1) % len(mp3Player.music)

	case kRepeatMode:
	}
	music := mp3Player.music[mp3Player.playIndex]
	mp3Player.appendToHistory(music)
	mp3Player.loadMusic(music)
}

func (c mp3Progress) execute(mp3Player *MP3Player) {
	bytes := int64(float64(mp3Player.mp3Size) * c.percent)
	bytes -= bytes % 4
	mp3Player.player.Seek(bytes, io.SeekStart)
	mp3Player.player.Play()
}

func (c m3pVolume) execute(mp3Player *MP3Player) {
	mp3Player.player.SetVolume(c.percent * c.percent)
}
