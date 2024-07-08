package cwidget

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type progressSlider struct {
	widget.BaseWidget
	slider *widget.Slider
	label  *widget.Label
}

func newProgressSlider() *progressSlider {
	s := progressSlider{
		slider: widget.NewSlider(0.0, 1.0),
		label:  widget.NewLabel("00:00"),
	}
	s.slider.Step = 0.01
	s.ExtendBaseWidget(&s)
	return &s
}

func (s *progressSlider) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, nil, s.label, s.slider))
}

func (s *progressSlider) setProgress(length time.Duration, percent float64) {
	mins := length / time.Minute
	secs := (length - mins*time.Minute) / time.Second
	s.label.SetText(fmt.Sprintf("%02d:%02d", mins, secs))
	s.slider.SetValue(percent)
}

func (s *progressSlider) setOnChanged(onChanged func(percent float64)) {
	s.slider.OnChangeEnded = onChanged
}
