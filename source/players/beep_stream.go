package players

import (
	"io"
	"log/slog"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

var _ beep.StreamCloser = &BeepStream{}

type BeepStream struct {
	stream     beep.StreamSeekCloser
	volumeCtrl *effects.Volume
	playCtrl   *beep.Ctrl
}

func newBeepStream(content io.ReadCloser, targetSampleRate beep.SampleRate, volume float64) *BeepStream {
	stream, format, err := mp3.Decode(content)
	if err != nil {
		content.Close()
		slog.Error("failed to create new beep stream", "error", err)
		return nil
	}

	resampled := beep.Resample(4, format.SampleRate, targetSampleRate, stream)
	vol := &effects.Volume{Streamer: resampled, Base: 2}
	ctrl := &beep.Ctrl{Streamer: vol}

	beepStream := BeepStream{
		stream:     stream,
		volumeCtrl: vol,
		playCtrl:   ctrl,
	}
	beepStream.setVolumeUnsafe(volume)
	return &beepStream
}

func (s *BeepStream) progress() float64 {
	speaker.Lock()
	defer speaker.Unlock()
	return float64(s.stream.Position()) / float64(s.stream.Len())
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
	percent = max(0.0, min(1.0, percent)) // Clamp input

	speaker.Lock()
	defer speaker.Unlock()

	targetSample := int(float64(s.stream.Len()) * percent)
	s.stream.Seek(targetSample)
	s.playCtrl.Paused = false
}

func (s *BeepStream) setVolume(percent float64) {
	speaker.Lock()
	defer speaker.Unlock()
	s.setVolumeUnsafe(percent)
}

func (s *BeepStream) setVolumeUnsafe(percent float64) {
	percent = max(0.0, min(1.0, percent))
	if percent <= 0.0 {
		s.volumeCtrl.Silent = true
		return
	}

	s.volumeCtrl.Silent = false
	if percent <= 0.7 {
		s.volumeCtrl.Volume = (percent / 0.7 * 6.0) - 6.0
	} else {
		s.volumeCtrl.Volume = ((percent - 0.7) / 0.3) * 1.5
	}
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
