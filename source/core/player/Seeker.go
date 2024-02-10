package player

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"meowyplayer.com/core/resource"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type Seeker struct {
	*mp3.Decoder
	oto.Player
}

func MakeSeeker(context *oto.Context, music *resource.Music) (Seeker, error) {
	mp3Data, err := os.ReadFile(resource.MusicPath(music))
	if err != nil {
		return Seeker{}, fmt.Errorf("can't read the music file %v", music.Title)
	}
	mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	if err != nil {
		return Seeker{}, fmt.Errorf("can't decode the music file %v", music.Title)
	}
	return Seeker{Decoder: mp3Decoder, Player: context.NewPlayer(mp3Decoder)}, nil
}

func (s *Seeker) CurrentProgressBytes() int64 {
	bytes, _ := s.Seek(0, io.SeekCurrent)
	return bytes
}

func (s *Seeker) CurrentProgressPercent() float64 {
	return float64(s.CurrentProgressBytes()) / float64(s.Length())
}

func (s *Seeker) SetProgress(percent float64) {
	bytes := int64(float64(s.Length()) * float64(percent))
	bytes -= bytes % 4
	s.Seek(bytes, io.SeekStart)
}

func (s *Seeker) IsOver() bool {
	return s.CurrentProgressBytes() == s.Length()
}

func (s *Seeker) PlayOrPause() {
	if s.IsPlaying() {
		s.Pause()
	} else {
		s.Play()
	}
}
