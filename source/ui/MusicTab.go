package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/manager"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cbinding"
	"meowyplayer.com/source/ui/cwidget"
)

func newMusicTab() *container.TabItem {
	data := cbinding.MakeMusicDataList()
	manager.GetCurrentAlbum().Attach(&data)

	return container.NewTabItemWithIcon("Music", resource.MusicTabIcon(), container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, cwidget.NewMusicTitleButton(&data, "Title"), cwidget.NewMusicDateButton(&data, "Date")),
			nil,
			container.NewGridWithRows(1, newMusicAdderLocalButton(&data), newMusicAdderOnlineButton(&data)),
			cwidget.NewMusicSearchBar(&data),
		),
		nil,
		nil,
		nil,
		newMusicViewList(&data),
	))
}

func newMusicViewList(data *cbinding.MusicDataList) *cwidget.MusicViewList {
	return cwidget.NewMusicViewList(data, func(m *player.Music) fyne.CanvasObject {
		view := cwidget.NewMusicView(m)
		view.OnTapped = func(*fyne.PointEvent) { manager.GetCurrentPlay().Set(player.NewPlay(data.GetAlbum(), m)) }
		view.OnTappedSecondary = func(*fyne.PointEvent) { showDeleteMusicDialog(m) }
		return view
	})
}

func newMusicAdderLocalButton(data *cbinding.MusicDataList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.MusicAdderLocalIcon(), func() {
		fileReader := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorIfAny(err)
			} else if result != nil {
				log.Printf("add %v from local to %v\n", result.URI().Name(), manager.GetCurrentAlbum().Get().Title)
				showErrorIfAny(manager.AddMusic(result))
			}
		}, getMainWindow())
		fileReader.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
		fileReader.SetConfirmText("Add")
		fileReader.Show()
	})
	button.Importance = widget.LowImportance
	return button
}

func newMusicAdderOnlineButton(data *cbinding.MusicDataList) *widget.Button {
	button := widget.NewButtonWithIcon("", resource.MusicAdderOnlineIcon(), func() {
		//to do
	})
	button.Importance = widget.LowImportance
	return button
}

func showDeleteMusicDialog(music *player.Music) {
	dialog.ShowConfirm("", fmt.Sprintf("Do you want to delete %v?", music.Title), func(delete bool) {
		if delete {
			log.Printf("delete %vfrom the album %v \n", music.Title, manager.GetCurrentAlbum().Get().Title)
			showErrorIfAny(manager.DeleteMusic(music))
		}
	}, getMainWindow())
}
