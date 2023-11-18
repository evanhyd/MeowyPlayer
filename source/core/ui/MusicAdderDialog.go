package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/network/fileformat"
	"meowyplayer.com/utility/network/scraper"
	"meowyplayer.com/utility/pattern"
)

func showAddLocalMusicDialog() {
	fileReader := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			showErrorIfAny(err)
		} else if result != nil {
			log.Printf("add %v from local to %v\n", result.URI().Name(), client.Manager().Album().Title)
			showErrorIfAny(client.AddMusicFromURIReader(client.Manager().Album(), result))
		}
	}, getWindow())
	fileReader.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	fileReader.SetConfirmText("Add")
	fileReader.Show()
}

func showAddOnlineMusicDialog() {
	//scraper menu
	var videoScraper scraper.VideoScraper

	//platform selector
	platformMenu := cwidget.NewDropDown("", resource.DefaultIcon)
	platformMenu.Add("YouTube", resource.YouTubeIcon, func() { videoScraper = scraper.NewClipzagScraper() })
	platformMenu.Add("BiliBili", resource.BiliBiliIcon, func() { fmt.Println("not implemented...") })
	platformMenu.Select(0)

	//video result view list
	videoResultData := pattern.Data[[]fileformat.VideoResult]{}
	videoResultViewList := cwidget.NewViewList(&videoResultData, container.NewVBox(),
		func(video fileformat.VideoResult) fyne.CanvasObject {
			return cwidget.NewVideoResultView(&video, fyne.NewSize(207, 128), func() {
				showErrorIfAny(client.DownloadMusic(client.Manager().Album(), &video))
			})
		},
	)

	//search bar
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search Video")
	searchEntry.ActionItem = cwidget.NewButtonWithIcon("", theme.SearchIcon(), func() { searchEntry.OnSubmitted(searchEntry.Text) })
	searchEntry.OnSubmitted = func(title string) {
		videos, err := videoScraper.Search(title)
		if err != nil {
			showErrorIfAny(err)
			return
		}
		videoResultData.Set(videos)
	}

	onlineMusicDialog := dialog.NewCustom("", "X", container.NewBorder(
		container.NewBorder(nil, nil, platformMenu, nil, searchEntry),
		nil,
		nil,
		nil,
		videoResultViewList,
	), getWindow())
	onlineMusicDialog.Resize(getWindow().Canvas().Size())
	onlineMusicDialog.Show()
}
