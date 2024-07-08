package cwidget

import (
	"playground/model"
	"playground/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MediaController struct {
	widget.BaseWidget
	preview        *AlbumCard
	title          *widget.RichText
	progressSlider *progressSlider
	volumeSlider   *volumeSlider
	prevButton     *widget.Button
	playButton     *widget.Button
	nextButton     *widget.Button
	modeButton     *widget.Button
}

func NewMediaController(media player.MediaPlayer) *MediaController {
	c := MediaController{
		preview:        NewAlbumCardConstructor(func(model.AlbumKey) {}, func(*fyne.PointEvent, model.AlbumKey) {})(),
		title:          widget.NewRichText(),
		progressSlider: newProgressSlider(),
		volumeSlider:   newVolumeSlider(),
		prevButton:     NewButton("", theme.MediaSkipPreviousIcon(), media.Prev),
		playButton:     NewButton("", theme.RadioButtonCheckedIcon(), media.Play),
		nextButton:     NewButton("", theme.MediaSkipNextIcon(), media.Next),
		modeButton:     NewButton("", theme.BrokenImageIcon(), nil),
	}
	c.volumeSlider.setOnChanged(media.SetVolume)
	c.progressSlider.setOnChanged(media.SetProgress)
	c.ExtendBaseWidget(&c)
	return &c
}

func (c *MediaController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil,
		nil,
		c.preview,
		nil,
		container.NewBorder(
			c.title,
			container.NewGridWithRows(1, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), c.modeButton, c.prevButton, c.playButton, c.nextButton), layout.NewSpacer(), c.volumeSlider),
			nil,
			nil,
			c.progressSlider,
		),
	))
}
