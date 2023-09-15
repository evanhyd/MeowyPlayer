package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MusicController struct {
	widget.BaseWidget

	title *widget.Label
	*progressController
	*buttonController
	*volumeController
}

func NewMusicController() *MusicController {
	controller := &MusicController{widget.BaseWidget{}, widget.NewLabel(""), newProgressController(), newButtonController(), newVolumeController()}
	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *MusicController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(
			c.title,
			container.NewGridWithRows(1, layout.NewSpacer(), c.buttonController, layout.NewSpacer(), c.volumeController),
			nil,
			nil,
			c.progressController,
		),
	)
}
