package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/oto/v2"
)

type volumeController struct {
	widget.BaseWidget
	volumeSlider *widget.Slider
	muteButton   *widget.Button
}

func newVolumeController() *volumeController {
	//volume slider
	slider := widget.NewSlider(0.0, 1.0)
	slider.Step = 0.01
	slider.Value = 0.5

	//mute button
	volume := 0.0
	button := widget.NewButtonWithIcon("", theme.MediaMusicIcon(), func() {
		if slider.Value == 0.0 {
			slider.SetValue(volume)
		} else {
			volume = slider.Value
			slider.SetValue(0.0)
		}
	})
	button.Importance = widget.LowImportance

	controller := &volumeController{widget.BaseWidget{}, slider, button}
	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *volumeController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, c.muteButton, nil, c.volumeSlider))
}

func (c *volumeController) BindVolume(mp3Player oto.Player) {
	mp3Player.SetVolume(c.volumeSlider.Value)
	c.volumeSlider.OnChanged = func(volume float64) { mp3Player.SetVolume(volume) }
}
