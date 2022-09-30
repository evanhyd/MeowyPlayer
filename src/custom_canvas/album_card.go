package custom_canvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/src/resource"
)

type AlbumCard struct {
	fyne.Container
}

func NewAlbumCard(title, description string, cardSize fyne.Size, primaryTap, secondaryTap func()) (*AlbumCard, error) {

	card := AlbumCard{}

	//load album card icon
	cardIcon, err := NewClickableIcon(canvas.NewImageFromFile(resource.GetAlbumIconPath(title)), cardSize, primaryTap, secondaryTap)
	if err != nil {
		return nil, err
	}

	//load album card title
	cardTitle := widget.NewLabel(title + "\n" + description)

	card.Container = *container.NewBorder(nil, nil, cardIcon, nil, cardTitle)
	return &card, nil
}
