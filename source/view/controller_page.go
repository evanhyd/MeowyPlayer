package view

import (
	"fmt"
	"image/color"
	"meowyplayer/model"
	"meowyplayer/player"
	"meowyplayer/view/internal/cwidget"
	"meowyplayer/view/internal/resource"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ControllerPage struct {
	widget.BaseWidget
	preview        *cwidget.AlbumCard
	title          *widget.RichText
	progressSlider *cwidget.ProgressSlider
	durationLabel  *widget.Label
	modeButton     *cwidget.DropDown
	prevButton     *widget.Button
	playButton     *widget.Button
	nextButton     *widget.Button
	volumeSlider   *cwidget.VolumeSlider

	music model.Music
}

func newControllerPage() *ControllerPage {
	var p ControllerPage
	p = ControllerPage{
		preview:        cwidget.NewAlbumCardConstructor(p.jumpToAlbum, func(*fyne.PointEvent, model.AlbumKey) {})(),
		title:          widget.NewRichText(),
		progressSlider: cwidget.NewProgressSlider(player.Instance().SetProgress),
		durationLabel:  widget.NewLabel(""),
		modeButton:     cwidget.NewDropDown(),
		prevButton:     cwidget.NewButton("", theme.MediaSkipPreviousIcon(), player.Instance().Prev),
		playButton:     cwidget.NewButton("", theme.RadioButtonCheckedIcon(), player.Instance().Play),
		nextButton:     cwidget.NewButton("", theme.MediaSkipNextIcon(), player.Instance().Next),
		volumeSlider:   cwidget.NewVolumeSlider(player.Instance().SetVolume),
	}

	p.modeButton.Add(cwidget.NewMenuItem(resource.RandomText(), resource.RandomIcon(), func() { player.Instance().SetMode(player.KRandomMode) }))
	p.modeButton.Add(cwidget.NewMenuItem(resource.OrderedText(), theme.MailForwardIcon(), func() { player.Instance().SetMode(player.KOrderedMode) }))
	p.modeButton.Add(cwidget.NewMenuItem(resource.RepeatText(), theme.MediaReplayIcon(), func() { player.Instance().SetMode(player.KRepeatMode) }))

	player.Instance().OnAlbumPlayed().AttachFunc(p.setCover)
	player.Instance().OnMusicPlayed().AttachFunc(p.setTitle)
	player.Instance().OnProgressUpdated().AttachFunc(p.setProgress)
	p.ExtendBaseWidget(&p)
	return &p
}

func (p *ControllerPage) CreateRenderer() fyne.WidgetRenderer {
	p.preview.HideTitle()
	p.title.Truncation = fyne.TextTruncateEllipsis
	frame := canvas.NewRectangle(color.Transparent)
	frame.SetMinSize(resource.KAlbumPreviewSize)

	return widget.NewSimpleRenderer(container.NewBorder(
		nil, nil,
		container.NewStack(frame, p.preview), nil,
		container.NewBorder(
			p.title,
			container.NewGridWithRows(1, layout.NewSpacer(), container.NewHBox(p.modeButton, p.prevButton, p.playButton, p.nextButton), layout.NewSpacer(), p.volumeSlider),
			nil,
			nil,
			container.NewBorder(nil, nil, nil, p.durationLabel, p.progressSlider),
		),
	))
}

func (p *ControllerPage) setCover(key model.AlbumKey) {
	album, err := model.UIClient().GetAlbum(key)
	if err != nil {
		//fail if the user deletes the album or switch the account
		fyne.LogError("failed to get album", err)
	}
	p.preview.Notify(album)
}

func (p *ControllerPage) setTitle(music model.Music) {
	p.title.Segments = p.title.Segments[:0]
	p.title.Segments = append(p.title.Segments, &widget.TextSegment{Text: music.Title(), Style: widget.RichTextStyleSubHeading})
	p.music = music
	p.Refresh()
}

func (p *ControllerPage) setProgress(percent float64) {
	remainDuration := p.music.Length() - time.Duration(float64(p.music.Length())*percent)
	mins := remainDuration / time.Minute
	secs := (remainDuration - mins*time.Minute) / time.Second
	p.durationLabel.SetText(fmt.Sprintf("%02d:%02d", mins, secs))
	p.progressSlider.SetValue(percent)
}

func (p *ControllerPage) jumpToAlbum(key model.AlbumKey) {
	if !key.IsEmpty() {
		if err := model.UIClient().SelectAlbum(key); err != nil {
			fyne.LogError("failed to select album", err)
		}
	}
}
