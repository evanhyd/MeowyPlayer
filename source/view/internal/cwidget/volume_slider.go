package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type VolumeSlider struct {
	widget.BaseWidget
	muteButton *widget.Button
	slider     *widget.Slider
}

func NewVolumeSlider(onChanged func(percent float64)) *VolumeSlider {
	s := VolumeSlider{
		muteButton: NewButton("", theme.VolumeUpIcon(), nil),
		slider:     widget.NewSlider(0.0, 1.0),
	}
	s.slider.Step = 0.01
	s.slider.Value = 0.5
	s.slider.OnChanged = func(percent float64) {
		onChanged(percent)
		if s.slider.Value == 0.0 {
			s.muteButton.SetIcon(theme.VolumeMuteIcon())
		} else {
			s.muteButton.SetIcon(theme.VolumeUpIcon())
		}
	}

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

func (v *VolumeSlider) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, v.muteButton, nil, v.slider))
}

func (v *VolumeSlider) SetVolume(volume float64) {
	v.slider.SetValue(volume)
}
