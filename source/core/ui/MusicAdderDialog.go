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
	"meowyplayer.com/utility/network/downloader"
	"meowyplayer.com/utility/network/fileformat"
	"meowyplayer.com/utility/network/scraper"
	"meowyplayer.com/utility/pattern"
)

func showAddLocalMusicDialog() {
	fileReader := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			showErrorIfAny(err)
		} else if result != nil {
			log.Printf("add %v from local to %v\n", result.URI().Name(), client.GetInstance().GetAlbum().Title)
			showErrorIfAny(client.AddMusicFromURIReader(result))
		}
	}, getWindow())
	fileReader.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	fileReader.SetConfirmText("Add")
	fileReader.Show()
}

func newVideoResultViewList(dataSource pattern.Subject[[]fileformat.VideoResult], onDownload func(videoResult *fileformat.VideoResult)) *cwidget.ViewList[fileformat.VideoResult] {
	return cwidget.NewViewList[fileformat.VideoResult](dataSource, container.NewVBox(),
		func(result fileformat.VideoResult) fyne.CanvasObject {
			return cwidget.NewVideoResultView(&result, fyne.NewSize(128.0*1.61803398875, 128.0), onDownload)
		},
	)
}

func showAddOnlineMusicDialog() {
	//scraper menu
	var videoScraper scraper.VideoScraper
	var musicDownloader downloader.MusicDownloader

	//video result data list
	videoResultData := pattern.Data[[]fileformat.VideoResult]{}
	videoResultViewList := newVideoResultViewList(&videoResultData,
		func(videoResult *fileformat.VideoResult) {
			musicData, err := musicDownloader.Download(videoResult)
			if err != nil {
				showErrorIfAny(err)
				return
			}
			showErrorIfAny(client.AddMusicFromDownloader(videoResult, musicData))
		},
	)

	platformMenu := cwidget.NewDropDown("", resource.DefaultIcon)
	platformMenu.Add("YouTube", resource.YouTubeIcon, func() {
		videoScraper = scraper.NewClipzagScraper()
		musicDownloader = downloader.NewY2MateDownloader()
	})
	platformMenu.Add("BiliBili", resource.BiliBiliIcon, func() {
		fmt.Println("not implemented...")
	})
	platformMenu.Select(0)

	//search bar
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Video")
	searchBar.ActionItem = cwidget.NewButtonWithIcon("", theme.SearchIcon(), func() { searchBar.OnSubmitted(searchBar.Text) })
	searchBar.OnSubmitted = func(title string) {
		result, err := videoScraper.Search(title)
		if err != nil {
			showErrorIfAny(err)
			return
		}
		videoResultData.Set(result)
	}

	onlineMusicDialog := dialog.NewCustom("", "X", container.NewBorder(
		container.NewBorder(nil, nil, platformMenu, nil, searchBar),
		nil,
		nil,
		nil,
		videoResultViewList,
	), getWindow())
	onlineMusicDialog.Resize(getWindow().Canvas().Size())
	onlineMusicDialog.Show()
}
