package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type volumeSlider struct {
	widget.BaseWidget
	slider     *widget.Slider
	muteButton *widget.Button
}

func newVolumeSlider() *volumeSlider {
	slider := widget.NewSlider(0.0, 1.0)
	slider.Step = 0.01
	slider.Value = 0.5

	volume := 0.0
	button := widget.NewButtonWithIcon("", theme.MediaMusicIcon(), nil)
	button.Importance = widget.LowImportance
	button.OnTapped = func() {
		if slider.Value == 0.0 {
			slider.SetValue(volume)
		} else {
			volume = slider.Value
			slider.SetValue(0.0)
		}
	}

	volumeSlider := &volumeSlider{slider: slider, muteButton: button}
	volumeSlider.ExtendBaseWidget(volumeSlider)
	return volumeSlider
}

func (v *volumeSlider) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, v.muteButton, nil, v.slider))
}

func (v *volumeSlider) SetOnChanged(onChanged func(volume float64)) {
	v.slider.OnChanged = onChanged
}

func (v *volumeSlider) Volume() float64 {
	return v.slider.Value
}
