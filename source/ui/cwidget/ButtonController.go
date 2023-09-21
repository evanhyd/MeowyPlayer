package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/source/resource"
)

type buttonController struct {
	widget.BaseWidget
	previousButton *widget.Button
	playButton     *widget.Button
	nextButton     *widget.Button
	modeButton     *widget.Button

	pauseManually bool
}

func newButtonController() *buttonController {
	controller := &buttonController{
		widget.BaseWidget{},
		widget.NewButton("<<", nil),
		widget.NewButton("O", nil),
		widget.NewButton(">>", nil),
		widget.NewButtonWithIcon("", resource.DefaultIcon(), nil),
		false,
	}
	controller.previousButton.Importance = widget.LowImportance
	controller.playButton.Importance = widget.LowImportance
	controller.nextButton.Importance = widget.LowImportance
	controller.modeButton.Importance = widget.LowImportance

	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *buttonController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(layout.NewSpacer(), c.modeButton, c.previousButton, c.playButton, c.nextButton))
}

func (c *buttonController) BindButton(mp3Decoder *mp3.Decoder, mp3Player oto.Player) {
	c.pauseManually = false

	c.playButton.OnTapped = func() {
		if c.pauseManually = mp3Player.IsPlaying(); c.pauseManually {
			mp3Player.Pause()
		} else {
			mp3Player.Play()
		}
	}

	c.nextButton.OnTapped = func() {
		c.pauseManually = false
		mp3Player.Pause()
	}
}

func (c *buttonController) IsPausedManually() bool {
	return c.pauseManually
}
