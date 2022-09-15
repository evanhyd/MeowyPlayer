package custom_canvas

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ClickableIcon struct {
	widget.Card
	primaryCallback   func()
	secondaryCallback func()
}

func NewClickableIcon(imagePath string, imageSize fyne.Size, onPrimaryTap, onSecondaryTap func()) (*ClickableIcon, error) {

	icon := ClickableIcon{}
	icon.ExtendBaseWidget(&icon)

	//load image
	icon.SetImage(canvas.NewImageFromFile(imagePath))
	icon.Image.SetMinSize(imageSize)

	//save callback
	icon.primaryCallback = onPrimaryTap
	icon.secondaryCallback = onSecondaryTap

	return &icon, nil
}

func (icon *ClickableIcon) Tapped(_ *fyne.PointEvent) {
	if icon.primaryCallback != nil {
		icon.primaryCallback()
	}
}

func (icon *ClickableIcon) TappedSecondary(_ *fyne.PointEvent) {
	if icon.secondaryCallback != nil {
		icon.secondaryCallback()
	}
}

func (icon *ClickableIcon) MouseIn(_ *desktop.MouseEvent) {
	icon.Image.Translucency = 0.3
	icon.Refresh()
}

func (icon *ClickableIcon) MouseMoved(*desktop.MouseEvent) {
}

func (icon *ClickableIcon) MouseOut() {
	icon.Image.Translucency = 0.0
	icon.Refresh()
}
