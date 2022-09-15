package custom_canvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/src/resource"
)

type AlbumCard struct {
	fyne.Container
}

func NewAlbumCard(title, description string, imageSize fyne.Size, primaryTap, secondaryTap func()) (*AlbumCard, error) {

	card := AlbumCard{}

	//load album card icon
	icon, err := NewClickableIcon(resource.GetAlbumIconPath(title), imageSize, primaryTap, secondaryTap)
	if err != nil {
		return nil, err
	}

	//load album card title
	cardTitle := widget.NewLabel(title + "\n" + description)

	card.Container = *container.NewBorder(nil, nil, icon, nil, cardTitle)
	return &card, nil
}
