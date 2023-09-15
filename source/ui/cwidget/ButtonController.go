package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/resource"
)

type buttonController struct {
	widget.BaseWidget
	previous   *widget.Button
	playButton *widget.Button
	nextButton *widget.Button
	modeButton *widget.Button
}

func newButtonController() *buttonController {
	controller := &buttonController{widget.BaseWidget{}, widget.NewButton("<<", nil), widget.NewButton("O", nil), widget.NewButton(">>", nil), widget.NewButtonWithIcon("", resource.DefaultIcon(), nil)}
	controller.previous.Importance = widget.LowImportance
	controller.playButton.Importance = widget.LowImportance
	controller.nextButton.Importance = widget.LowImportance
	controller.modeButton.Importance = widget.LowImportance

	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *buttonController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(layout.NewSpacer(), c.modeButton, c.previous, c.playButton, c.nextButton))
}
