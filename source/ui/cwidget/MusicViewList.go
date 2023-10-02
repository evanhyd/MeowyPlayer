package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
)

type MusicViewList struct {
	widget.BaseWidget
	box      *fyne.Container
	makeView func(*resource.Music) fyne.CanvasObject
}

func NewMusicViewList(data *cbinding.MusicDataList, makeMusicView func(*resource.Music) fyne.CanvasObject) *MusicViewList {
	list := &MusicViewList{box: container.NewVBox(), makeView: makeMusicView}
	data.Attach(list)
	list.ExtendBaseWidget(list)
	return list
}

func (m *MusicViewList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewScroll(m.box))
}

func (m *MusicViewList) Notify(data []resource.Music) {
	m.box.RemoveAll()
	views := m.makeViews(data)
	for i := range views {
		m.box.Add(views[i])
	}
}

func (m *MusicViewList) makeViews(data []resource.Music) []fyne.CanvasObject {
	views := make([]fyne.CanvasObject, len(data))
	for i := range data {
		views[i] = m.makeView(&data[i])
	}
	return views
}
