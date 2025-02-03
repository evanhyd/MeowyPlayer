package player

import (
	"io"
	"math"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
)

type ControllableMp3Stream struct {
	resource   beep.StreamSeekCloser
	volumeCtrl *effects.Volume
	playCtrl   *beep.Ctrl
}

func newMp3Stream(rsc io.ReadSeekCloser) (*ControllableMp3Stream, error) {
	resource, format, err := mp3.Decode(rsc)
	if err != nil {
		return nil, err
	}
	resamplerCtrl := beep.Resample(16, format.SampleRate, kSampleRate, resource)
	volumeCtrl := &effects.Volume{Streamer: resamplerCtrl, Base: 2}
	playCtrl := &beep.Ctrl{Streamer: volumeCtrl}
	return &ControllableMp3Stream{resource, volumeCtrl, playCtrl}, nil
}

func (s *ControllableMp3Stream) setVolume(percent float64) {
	const kVolumeOffsetPercent = -0.7
	fixedPercent := percent + kVolumeOffsetPercent
	s.volumeCtrl.Volume = 10 * math.Copysign(fixedPercent*fixedPercent, fixedPercent)
	s.volumeCtrl.Silent = (percent == 0.0)
}

func (s *ControllableMp3Stream) Stream(sample [][2]float64) (int, bool) {
	return s.playCtrl.Stream(sample)
}

func (s *ControllableMp3Stream) Close() error {
	return s.resource.Close()
}

func (s *ControllableMp3Stream) Err() error {
	return s.playCtrl.Err()
}
