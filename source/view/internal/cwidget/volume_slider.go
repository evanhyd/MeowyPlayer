package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type volumeSlider struct {
	widget.BaseWidget
	muteButton *widget.Button
	slider     *widget.Slider
}

func newVolumeSlider() *volumeSlider {
	s := volumeSlider{
		muteButton: NewButton("", theme.VolumeUpIcon(), nil),
		slider:     widget.NewSlider(0.0, 1.0),
	}
	s.slider.Step = 0.01
	s.slider.Value = 0.5
	prevPercent := 0.0
	s.muteButton.OnTapped = func() {
		if s.slider.Value == 0.0 {
			s.slider.SetValue(prevPercent)
		} else {
			prevPercent = s.slider.Value
			s.slider.SetValue(0.0)
		}
	}

	s.ExtendBaseWidget(&s)
	return &s
}

func (v *volumeSlider) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, v.muteButton, nil, v.slider))
}

func (v *volumeSlider) setOnChanged(onChanged func(percent float64)) {
	v.slider.OnChanged = func(percent float64) {
		onChanged(percent)
		if v.slider.Value == 0.0 {
			v.muteButton.SetIcon(theme.VolumeMuteIcon())
		} else {
			v.muteButton.SetIcon(theme.VolumeUpIcon())
		}
	}
}

func (v *volumeSlider) setVolume(volume float64) {
	v.slider.SetValue(volume)
}
