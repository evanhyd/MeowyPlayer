package custom_canvas

import (
	"fmt"

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

	//play button
	playIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("music_card_play.png"))
	if err != nil {
		return nil, err
	}
	play := widget.NewButtonWithIcon("", playIcon, playTapped)

	//delete button
	deleteIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("music_card_delete.png"))
	if err != nil {
		return nil, err
	}
	delete := widget.NewButtonWithIcon("", deleteIcon, deleteTapped)

	//music duraiton
	musicDuration := widget.NewLabel(fmt.Sprintf("%02d:%02d", duration/60, duration%60))

	//music title
	musicTitle := widget.NewLabel(title)

	card.Container = *container.NewHBox(play, delete, musicDuration, musicTitle)
	return &card, nil
}
