package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ProgressSlider struct {
	widget.Slider
	onReleased        func(float64)
	isBeingControlled bool
}

func NewProgressSlider(onReleased func(percent float64)) *ProgressSlider {
	s := ProgressSlider{Slider: widget.Slider{Min: 0.0, Max: 1.0, Step: 0.001}, onReleased: onReleased}
	s.ExtendBaseWidget(&s)
	return &s
}

func (s *ProgressSlider) Tapped(e *fyne.PointEvent) {
	s.Slider.Tapped(e)
	s.onReleased(s.Value)
}

func (s *ProgressSlider) Dragged(e *fyne.DragEvent) {
	s.isBeingControlled = true
	s.Slider.Dragged(e)
}

func (s *ProgressSlider) DragEnd() {
	if s.isBeingControlled {
		s.onReleased(s.Value)
		s.isBeingControlled = false
	}
}

func (s *ProgressSlider) setValue(percent float64) {
	if !s.isBeingControlled {
		s.Slider.SetValue(percent)
	}
}
