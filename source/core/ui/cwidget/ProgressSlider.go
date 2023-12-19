package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ProgressSlider struct {
	widget.Slider
	OnReleased        func(float64)
	isBeingControlled bool
}

func NewProgressSlider(step float64) *ProgressSlider {
	progressSlider := &ProgressSlider{Slider: widget.Slider{Min: 0.0, Max: 1.0, Step: step}}
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
	s.isBeingControlled = true
	s.Slider.Dragged(e)
}

func (s *ProgressSlider) DragEnd() {
	if s.OnReleased != nil {
		s.OnReleased(s.Value)
	}
	s.isBeingControlled = false
}

func (s *ProgressSlider) SetValue(percent float64) {
	if !s.isBeingControlled {
		s.Slider.SetValue(percent)
	}
}
