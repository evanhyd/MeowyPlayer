package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
)

type MenuChannel struct {
	Progress chan float64
	Volume   chan float64
	Play     chan struct{}
	Rollback chan struct{}
	Skip     chan struct{}
	Mode     chan struct{}
}

func makeMenuChannel() MenuChannel {
	return MenuChannel{
		make(chan float64, 16),
		make(chan float64, 16),
		make(chan struct{}, 16),
		make(chan struct{}, 16),
		make(chan struct{}, 16),
		make(chan struct{}, 16),
	}
}

type PlayerMenu struct {
	widget.BaseWidget
	title          *widget.Label
	progressSlider *ProgressSlider
	durationLabel  *widget.Label
	rollbackButton *widget.Button
	playButton     *widget.Button
	skipButton     *widget.Button
	modeButton     *widget.Button
	volumeSlider   *volumeSlider
	controlChannel MenuChannel
}

func NewPlayerMenu() *PlayerMenu {
	menuChannel := makeMenuChannel()

	title := widget.NewLabel("")

	progressSlider := NewProgressSlider(0.0, 1.0, 0.001, 0.0)
	progressSlider.OnReleased = func(percent float64) { menuChannel.Progress <- percent }

	durationLabel := widget.NewLabel("00:00")

	rollbackButton := widget.NewButton("<<", func() { menuChannel.Rollback <- struct{}{} })
	rollbackButton.Importance = widget.LowImportance

	playButton := widget.NewButton("O", func() { menuChannel.Play <- struct{}{} })
	playButton.Importance = widget.LowImportance

	skipButton := widget.NewButton(">>", func() { menuChannel.Skip <- struct{}{} })
	skipButton.Importance = widget.LowImportance

	modeButton := widget.NewButtonWithIcon("", resource.DefaultIcon(), func() { menuChannel.Mode <- struct{}{} })
	modeButton.Importance = widget.LowImportance

	volumeSlider := newVolumeSlider()
	volumeSlider.SetOnChanged(func(volume float64) { menuChannel.Volume <- volume })

	menu := &PlayerMenu{
		title:          title,
		progressSlider: progressSlider,
		durationLabel:  durationLabel,
		rollbackButton: rollbackButton,
		playButton:     playButton,
		skipButton:     skipButton,
		modeButton:     modeButton,
		volumeSlider:   volumeSlider,
		controlChannel: menuChannel,
	}
	menu.ExtendBaseWidget(menu)
	return menu
}

func (c *PlayerMenu) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		c.title,
		container.NewGridWithRows(1, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), c.rollbackButton, c.playButton, c.skipButton, c.modeButton), layout.NewSpacer(), c.volumeSlider),
		nil,
		nil,
		container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider),
	))
}

func (c *PlayerMenu) GetMenuChannel() MenuChannel {
	return c.controlChannel
}

func (c *PlayerMenu) SetMusic(music *player.Music) {
	c.title.SetText(music.SimpleTitle())
	c.progressSlider.SetValue(0.0)
	c.controlChannel.Volume <- c.volumeSlider.Volume()
}

func (c *PlayerMenu) UpdateProgressBar(percent float64) {
	c.progressSlider.SetValue(percent)
}
