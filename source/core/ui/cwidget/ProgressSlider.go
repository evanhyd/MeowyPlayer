package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ProgressSlider struct {
	widget.Slider
	OnReleased func(float64)
	pressed    bool
}

func NewProgressSlider(min, max, step, value float64) *ProgressSlider {
	progressSlider := &ProgressSlider{Slider: widget.Slider{Value: value, Min: min, Max: max, Step: step}}
	progressSlider.ExtendBaseWidget(progressSlider)
	return progressSlider
}

func (s *ProgressSlider) Dragged(e *fyne.DragEvent) {
	s.pressed = true
	s.Slider.Dragged(e)
}

func (s *ProgressSlider) DragEnd() {
	if s.OnReleased != nil {
		s.OnReleased(s.Value)
	}
	s.pressed = false
}

func (s *ProgressSlider) SetValue(percent float64) {
	if !s.pressed {
		s.Slider.SetValue(percent)
	}
}
