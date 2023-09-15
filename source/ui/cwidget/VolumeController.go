package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type volumeController struct {
	widget.BaseWidget
	muteButton   *widget.Button
	volumeSlider *widget.Slider
}

func newVolumeController() *volumeController {
	controller := &volumeController{widget.BaseWidget{}, widget.NewButtonWithIcon("", theme.MediaMusicIcon(), nil), widget.NewSlider(0.0, 1.0)}
	controller.muteButton.Importance = widget.LowImportance
	controller.volumeSlider.Step = 0.01
	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *volumeController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, c.muteButton, nil, c.volumeSlider))
}
