package client

import (
	"io"
	"sync"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type MP3Controller struct {
	sync.Mutex
	*mp3.Decoder
	oto.Player
}

func makeMP3Player(decoder *mp3.Decoder, player oto.Player) MP3Controller {
	return MP3Controller{Decoder: decoder, Player: player}
}

func (m *MP3Controller) CurrentProgressBytes() int64 {
	bytes, _ := m.Seek(0, io.SeekCurrent)
	return bytes
}

func (m *MP3Controller) CurrentProgressPercent() float64 {
	bytes, _ := m.Seek(0, io.SeekCurrent)
	return float64(bytes) / float64(m.Length())
}

func (m *MP3Controller) SetProgress(percent float64) {
	m.Lock()
	defer m.Unlock()
	bytes := int64(float64(m.Length()) * float64(percent))
	bytes -= bytes % 4
	m.Seek(bytes, io.SeekStart)
}

func (m *MP3Controller) IsOver() bool {
	return m.CurrentProgressBytes() == m.Length()
}

func (m *MP3Controller) PlayOrPause() {
	if m.IsPlaying() {
		m.Pause()
	} else {
		m.Play()
	}
}
