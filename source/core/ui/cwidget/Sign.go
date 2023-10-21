package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Sign struct {
	widget.BaseWidget
	title *widget.Label
	icon  *widget.Icon
}

func NewSign(title string, icon fyne.Resource) *Sign {
	sign := &Sign{title: widget.NewLabel(title), icon: widget.NewIcon(icon)}
	sign.ExtendBaseWidget(sign)
	return sign
}

func (s *Sign) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, s.icon, nil, s.title))
}
