package player

import (
	"bytes"
	"io"
	"os"
	"sync"

	"meowyplayer.com/core/resource"
	"meowyplayer.com/utility/logger"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type MP3Controller struct {
	mutex sync.Mutex
	*mp3.Decoder
	oto.Player
}

func NewMP3Controller(context *oto.Context, music *resource.Music) *MP3Controller {
	mp3Data, err := os.ReadFile(resource.MusicPath(music))
	if err != nil {
		logger.Error(err, 0)
	}
	mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	if err != nil {
		logger.Error(err, 0)
	}
	return &MP3Controller{Decoder: mp3Decoder, Player: context.NewPlayer(mp3Decoder)}
}

func (m *MP3Controller) CurrentProgressBytes() int64 {
	bytes, _ := m.Seek(0, io.SeekCurrent)
	return bytes
}

func (m *MP3Controller) CurrentProgressPercent() float64 {
	return float64(m.CurrentProgressBytes()) / float64(m.Length())
}

func (m *MP3Controller) SetProgress(percent float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
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
