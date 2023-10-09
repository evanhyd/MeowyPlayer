package cwidget

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/resource"
)

type MediaCommander interface {
	CommandProgress(float64)
	CommandVolume(float64)
	CommandPlay()
	CommandRollback()
	CommandSkip()
	CommandMode(int)
}

type MediaMenu struct {
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

func NewMediaMenu() *MediaMenu {
	modeIcons := []fyne.Resource{resource.RandomIcon(), theme.MailForwardIcon(), theme.ViewRefreshIcon()}
	menu := &MediaMenu{
		widget.BaseWidget{},
		widget.NewLabel(""),
		NewProgressSlider(0.0, 1.0, 0.001, 0.0),
		widget.NewLabel("00:00"),
		newModeButton(nil, modeIcons, nil),
		NewButtonWithIcon("", theme.MediaFastRewindIcon(), nil),
		NewButtonWithIcon("", theme.RadioButtonCheckedIcon(), nil),
		NewButtonWithIcon("", theme.MediaFastForwardIcon(), nil),
		newVolumeSlider(),
	}
	menu.ExtendBaseWidget(menu)
	return menu
}

func (c *MediaMenu) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		c.title,
		container.NewGridWithRows(1, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), c.modeButton, c.rollbackButton, c.playButton, c.skipButton), layout.NewSpacer(), c.volumeSlider),
		nil,
		nil,
		container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider),
	))
}

func (c *MediaMenu) Bind(commander MediaCommander) {
	c.progressSlider.OnReleased = func(percent float64) { commander.CommandProgress(percent) }
	c.playButton.OnTapped = func() { commander.CommandPlay() }
	c.skipButton.OnTapped = func() { commander.CommandSkip() }
	c.rollbackButton.OnTapped = func() { commander.CommandRollback() }
	c.modeButton.OnTapped = func(mode int) { commander.CommandMode(mode) }
	c.volumeSlider.SetOnChanged(func(volume float64) { commander.CommandVolume(volume) })
}

func (c *MediaMenu) SetMusic(music *resource.Music) {
	c.title.SetText(music.SimpleTitle())
	c.progressSlider.SetValue(0.0)
}

func (c *MediaMenu) UpdateProgress(length time.Duration, percent float64) {
	const kConversionFactor = 60
	length = time.Duration(float64(length) * percent)
	mins := int(length.Minutes()) % kConversionFactor
	secs := int(length.Seconds()) % kConversionFactor
	c.durationLabel.SetText(fmt.Sprintf("%02v:%02v", mins, secs))
	c.progressSlider.SetValue(percent)
}

func (c *MediaMenu) Volume() float64 {
	return c.volumeSlider.Volume()
}
