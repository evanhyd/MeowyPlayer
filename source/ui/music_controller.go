package ui

import (
	"meowyplayer/ui/internal/mcontainer"
	"meowyplayer/ui/internal/mwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MusicController struct {
	widget.BaseWidget
	playlistCover  *canvas.Image
	title          *widget.RichText
	progressSlider *widget.Slider
	durationLabel  *widget.Label
	modeButton     *widget.Button
	skipPrevButton *widget.Button
	playButton     *widget.Button
	skipNextButton *widget.Button
	volumeSlider   *widget.Slider
}

func newMusicController() *MusicController {
	c := MusicController{
		playlistCover:  canvas.NewImageFromResource(resourceIconPng),
		title:          widget.NewRichTextWithText(""),
		progressSlider: widget.NewSlider(0.0, 1.0),
		durationLabel:  widget.NewLabel("00:00"),
		modeButton:     widget.NewButtonWithIcon("", theme.BrokenImageIcon(), nil),
		skipPrevButton: widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), nil),
		playButton:     widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		skipNextButton: widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), nil),
		volumeSlider:   widget.NewSlider(0.0, 1.0),
	}

	c.playlistCover.SetMinSize(mwidget.PlaylistCardSize)
	c.title.Truncation = fyne.TextTruncateEllipsis
	c.modeButton.Importance = widget.LowImportance
	c.skipPrevButton.Importance = widget.LowImportance
	c.playButton.Importance = widget.LowImportance
	c.skipNextButton.Importance = widget.LowImportance

	c.ExtendBaseWidget(&c)
	return &c
}

func (c *MusicController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil,
		nil,
		c.playlistCover,
		nil,
		container.NewBorder(
			c.title,
			mcontainer.NewHSplit(0.8, container.NewCenter(container.NewHBox(c.modeButton, c.skipPrevButton, c.playButton, c.skipNextButton)), c.volumeSlider),
			nil,
			nil,
			container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider),
		),
	))
}
