package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ProgressSlider struct {
	widget.Slider
	OnReleased    func(float64)
	isUserControl bool
}

func NewProgressSlider(min, max, step, value float64) *ProgressSlider {
	progressSlider := &ProgressSlider{Slider: widget.Slider{Value: value, Min: min, Max: max, Step: step}}
	progressSlider.ExtendBaseWidget(progressSlider)
	return progressSlider
}

func (s *ProgressSlider) Tapped(e *fyne.PointEvent) {
	s.Slider.Tapped(e)
	if s.OnReleased != nil {
		s.OnReleased(s.Value)
	}
}

func (s *ProgressSlider) Dragged(e *fyne.DragEvent) {
	s.isUserControl = true
	s.Slider.Dragged(e)
}

func (s *ProgressSlider) DragEnd() {
	if s.OnReleased != nil {
		s.OnReleased(s.Value)
	}
	s.isUserControl = false
}

func (s *ProgressSlider) SetValue(percent float64) {
	if !s.isUserControl {
		s.Slider.SetValue(percent)
	}
}
