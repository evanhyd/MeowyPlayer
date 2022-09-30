package custom_canvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type YoutubeCard struct {
	fyne.Container
}

func NewYoutubeCard(title, imageURL string, imageSize fyne.Size, primaryTap, secondaryTap func()) (*YoutubeCard, error) {

	card := &YoutubeCard{}

	//load thumbnail
	res, err := fyne.LoadResourceFromURLString(imageURL)
	if err != nil {
		return nil, err
	}
	cardIcon, err := NewClickableIcon(canvas.NewImageFromResource(res), imageSize, primaryTap, secondaryTap)
	if err != nil {
		return nil, err
	}

	//load youtube card title
	cardTitle := widget.NewLabel(title)

	card.Container = *container.NewBorder(nil, nil, cardIcon, nil, cardTitle)
	return card, nil
}
