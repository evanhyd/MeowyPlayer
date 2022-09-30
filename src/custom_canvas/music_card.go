package custom_canvas

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"meowyplayer.com/src/resource"
)

type MusicCard struct {
	fyne.Container
}

func NewMusicCard(title string, duration int64, playTapped, deleteTapped func()) (*MusicCard, error) {
	card := MusicCard{}

	//delete button
	deleteIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("music_card_delete.png"))
	if err != nil {
		return nil, err
	}
	delete := widget.NewButtonWithIcon("", deleteIcon, deleteTapped)

	//play button
	title = title[:strings.LastIndex(title, ".mp3")]
	play := widget.NewButton(fmt.Sprintf("%02d:%02d", duration/60, duration%60)+" | "+title, playTapped)
	play.Importance = widget.LowImportance
	play.Alignment = widget.ButtonAlignLeading

	card.Container = *container.NewBorder(nil, nil, delete, nil, play)
	return &card, nil
}
