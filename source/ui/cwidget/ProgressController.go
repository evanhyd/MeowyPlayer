package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type progressController struct {
	widget.BaseWidget
	progressSlider *widget.Slider
	durationLabel  *widget.Label
}

func newProgressController() *progressController {
	controller := &progressController{widget.BaseWidget{}, widget.NewSlider(0.0, 1.0), widget.NewLabel("00:00")}
	controller.progressSlider.Step = 0.001
	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *progressController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider))
}
