package cwidget

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/resource"
)

type MenuCommander interface {
	CommandProgress(float64)
	CommandVolume(float64)
	CommandPlay()
	CommandRollback()
	CommandSkip()
	CommandMode(int)
}

type CommandMenu struct {
	widget.BaseWidget
	title          *widget.Label
	progressSlider *ProgressSlider
	durationLabel  *widget.Label
	rollbackButton *widget.Button
	playButton     *widget.Button
	skipButton     *widget.Button
	modeButton     *ModeButton
	volumeSlider   *volumeSlider
}

func NewCommandMenu() *CommandMenu {
	modeIcons := []fyne.Resource{resource.PlayModeRandomIcon(), resource.PlayModeOrderedIcon(), resource.PlayModeRepeatIcon()}
	menu := &CommandMenu{
		widget.BaseWidget{},
		widget.NewLabel(""),
		NewProgressSlider(0.0, 1.0, 0.001, 0.0),
		widget.NewLabel("00:00"),
		NewButton("<<", nil),
		NewButton("O", nil),
		NewButton(">>", nil),
		newModeButton(nil, modeIcons, nil),
		newVolumeSlider(),
	}
	menu.ExtendBaseWidget(menu)
	return menu
}

func (c *CommandMenu) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		c.title,
		container.NewGridWithRows(1, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), c.rollbackButton, c.playButton, c.skipButton, c.modeButton), layout.NewSpacer(), c.volumeSlider),
		nil,
		nil,
		container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider),
	))
}

func (c *CommandMenu) Bind(commander MenuCommander) {
	c.progressSlider.OnReleased = func(percent float64) { commander.CommandProgress(percent) }
	c.playButton.OnTapped = func() { commander.CommandPlay() }
	c.skipButton.OnTapped = func() { commander.CommandSkip() }
	c.rollbackButton.OnTapped = func() { commander.CommandRollback() }
	c.modeButton.OnTapped = func(mode int) { commander.CommandMode(mode) }
	c.volumeSlider.SetOnChanged(func(volume float64) { commander.CommandVolume(volume) })
}

func (c *CommandMenu) SetMusic(music *resource.Music) {
	c.title.SetText(music.SimpleTitle())
	c.progressSlider.SetValue(0.0)
	c.volumeSlider.SetVolume(c.volumeSlider.Volume())
}

func (c *CommandMenu) UpdateProgress(length time.Duration, percent float64) {
	const kConversionFactor = 60
	length = time.Duration(float64(length) * percent)
	mins := int(length.Minutes()) % kConversionFactor
	secs := int(length.Seconds()) % kConversionFactor
	c.durationLabel.SetText(fmt.Sprintf("%02v:%02v", mins, secs))
	c.progressSlider.SetValue(percent)
}
