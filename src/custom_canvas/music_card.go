package custom_canvas

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type MusicCard struct {
	widget.Button
	primaryCallback   func()
	secondaryCallback func()
}

func NewMusicCard(title string, duration int64, onPrimaryTap, onSecondaryTap func()) (*MusicCard, error) {
	card := MusicCard{}
	card.ExtendBaseWidget(&card)

	title = title[:strings.LastIndex(title, ".mp3")]
	card.SetText(fmt.Sprintf("%02d:%02d", duration/60, duration%60) + " | " + title)
	card.Importance = widget.LowImportance
	card.Alignment = widget.ButtonAlignLeading
	card.primaryCallback = onPrimaryTap
	card.secondaryCallback = onSecondaryTap
	return &card, nil
}

func (card *MusicCard) Tapped(_ *fyne.PointEvent) {
	if card.primaryCallback != nil {
		card.primaryCallback()
	}
}

func (card *MusicCard) TappedSecondary(_ *fyne.PointEvent) {
	if card.secondaryCallback != nil {
		card.secondaryCallback()
	}
}
