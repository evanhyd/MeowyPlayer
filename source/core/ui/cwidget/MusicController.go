package cwidget

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/resource"
)

type MediaController interface {
	OnProgress(percent float64)
	OnVolume(volume float64)
	OnPlay()
	OnRollback()
	OnSkip()
	OnMode(mode int)
}

type MusicController struct {
	widget.BaseWidget
	title          *widget.Label
	progressSlider *ProgressSlider
	durationLabel  *widget.Label
	modeButton     *ModeButton
	rollbackButton *widget.Button
	playButton     *widget.Button
	skipButton     *widget.Button
	volumeSlider   *volumeSlider
}

func NewMusicController() *MusicController {
	modeIcons := []fyne.Resource{resource.RandomIcon, theme.MailForwardIcon(), theme.ViewRefreshIcon()}
	menu := &MusicController{
		widget.BaseWidget{},
		widget.NewLabel(""),
		NewProgressSlider(0.001),
		widget.NewLabel("00:00"),
		newModeButton(nil, modeIcons, nil),
		NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil),
		NewButtonWithIcon("", theme.RadioButtonCheckedIcon(), nil),
		NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil),
		newVolumeSlider(),
	}
	menu.ExtendBaseWidget(menu)
	return menu
}

func (c *MusicController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		c.title,
		container.NewGridWithRows(1, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), c.modeButton, c.rollbackButton, c.playButton, c.skipButton), layout.NewSpacer(), c.volumeSlider),
		nil,
		nil,
		container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider),
	))
}

func (c *MusicController) Bind(controller MediaController) {
	c.progressSlider.OnReleased = controller.OnProgress
	c.playButton.OnTapped = controller.OnPlay
	c.skipButton.OnTapped = controller.OnSkip
	c.rollbackButton.OnTapped = controller.OnRollback
	c.modeButton.OnTapped = controller.OnMode
	c.volumeSlider.SetOnChanged(controller.OnVolume)
}

func (c *MusicController) SetMusic(music *resource.Music) {
	c.title.SetText(music.SimpleTitle())
	c.progressSlider.SetValue(0.0)
}

func (c *MusicController) UpdateProgress(length time.Duration, percent float64) {
	const kConversionFactor = 60
	length = time.Duration(float64(length) * percent)
	mins := int(length.Minutes()) % kConversionFactor
	secs := int(length.Seconds()) % kConversionFactor
	c.durationLabel.SetText(fmt.Sprintf("%02v:%02v", mins, secs))
	c.progressSlider.SetValue(percent)
}

func (c *MusicController) Volume() float64 {
	return c.volumeSlider.Volume()
}
