package player

import (
	"io"
	"log"
	"math"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

var _ beep.StreamCloser = &BeepStream{}

type BeepStream struct {
	stream     beep.StreamSeekCloser
	volumeCtrl effects.Volume
	playCtrl   beep.Ctrl
}

func newBeepStream(content io.ReadSeekCloser, targetSampleRate beep.SampleRate) *BeepStream {
	stream, format, err := mp3.Decode(content)
	if err != nil {
		log.Fatalf("failed to create new beep stream: %v\n", err)
	}
	resampled := beep.Resample(16, format.SampleRate, targetSampleRate, stream)

	return &BeepStream{
		stream,
		effects.Volume{Streamer: resampled, Base: 2},
		beep.Ctrl{Streamer: resampled},
	}
}

func (s *BeepStream) resume() {
	speaker.Lock()
	defer speaker.Unlock()
	s.playCtrl.Paused = false
}

func (s *BeepStream) suspend() {
	speaker.Lock()
	defer speaker.Unlock()
	s.playCtrl.Paused = true
}

func (s *BeepStream) setProgress(percent float64) {
	speaker.Lock()
	defer speaker.Unlock()
	bytes := int(float64(s.stream.Len()) * percent) // Align to 4 bytes.
	bytes -= bytes % 4
	s.stream.Seek(bytes)
	s.playCtrl.Paused = false
}

func (s *BeepStream) setVolume(percent float64) {
	speaker.Lock()
	defer speaker.Unlock()
	const volumeOffset = -0.7
	fixedPercent := percent + volumeOffset
	s.volumeCtrl.Volume = 10 * math.Copysign(fixedPercent*fixedPercent, fixedPercent)
	s.volumeCtrl.Silent = (percent == 0.0)
}

func (s *BeepStream) Stream(sample [][2]float64) (int, bool) {
	return s.playCtrl.Stream(sample)
}

func (s *BeepStream) Close() error {
	return s.stream.Close()
}

func (s *BeepStream) Err() error {
	return s.playCtrl.Err()
}
