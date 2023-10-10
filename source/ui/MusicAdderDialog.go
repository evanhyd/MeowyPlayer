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
	"meowyplayer.com/source/client"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/ui/cwidget"
	"meowyplayer.com/utility/assert"
	"meowyplayer.com/utility/network/scraper"
	"meowyplayer.com/utility/pattern"
)

func showAddLocalMusicDialog() {
	fileReader := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			showErrorIfAny(err)
		} else if result != nil {
			log.Printf("add %v from local to %v\n", result.URI().Name(), client.GetAlbumData().Get().Title)
			showErrorIfAny(client.AddLocalMusic(result))
		}
	}, getWindow())
	fileReader.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	fileReader.SetConfirmText("Add")
	fileReader.Show()
}

func newVideoResultViewList(dataSource pattern.Subject[[]scraper.VideoResult]) *cwidget.ViewList[scraper.VideoResult] {
	return cwidget.NewViewList[scraper.VideoResult](dataSource, container.NewVBox(),
		func(result scraper.VideoResult) fyne.CanvasObject {
			return cwidget.NewVideoResultView(&result, fyne.NewSize(128.0*1.61803398875, 128.0))
		},
	)
}

func showAddOnlineMusicDialog() {
	//video result data list
	videoResultData := pattern.Data[[]scraper.VideoResult]{}
	videoResultViewList := newVideoResultViewList(&videoResultData)

	//scraper menu
	var videoScraper scraper.VideoScraper
	scraperMenu := cwidget.NewDropDown("", resource.DefaultIcon())
	scraperMenu.Add("YouTube", resource.YouTubeIcon(), func() { videoScraper = scraper.NewClipzagScraper() })
	scraperMenu.Add("BiliBili", resource.BiliBiliIcon(), func() { fmt.Println("not implemented...") })
	scraperMenu.Select(0)

	//search bar
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Video")
	searchBar.ActionItem = cwidget.NewButtonWithIcon("", theme.SearchIcon(), func() { searchBar.OnSubmitted(searchBar.Text) })
	searchBar.OnSubmitted = func(title string) {
		result, err := videoScraper.Search(title)
		assert.NoErr(err, "failed to scrape the video info")
		videoResultData.Set(result)
	}

	onlineMusicDialog := dialog.NewCustom("", "O", container.NewBorder(
		container.NewBorder(nil, nil, scraperMenu, nil, searchBar),
		nil,
		nil,
		nil,
		videoResultViewList,
	), getWindow())
	onlineMusicDialog.Resize(getWindow().Canvas().Size())
	onlineMusicDialog.Show()

	// 	func(result scraper.ClipzagResult, canvas fyne.CanvasObject) {
	// 		borderItems := canvas.(*fyne.Container).Objects
	// 		gridItems := borderItems[0].(*fyne.Container).Objects

	// 		videoTitle := gridItems[0].(*widget.Label)
	// 		if videoTitle.Text != result.VideoTitle() {
	// 			card := borderItems[1].(*widget.Card)
	// 			card.Image = result.Thumbnail()

	// 			videoTitle.Text = result.VideoTitle()

	// 			videoStats := gridItems[1].(*widget.Label)
	// 			videoStats.Text = result.Stats()

	// 			description := gridItems[2].(*widget.Label)
	// 			description.Text = result.Description()

	// 			canvas.Refresh()
	// 		}
	// 	},
	// )

	// scroll.SetOnSelected(func(result *scraper.ClipzagResult) {
	// 	progress := dialog.NewCustom(result.VideoTitle(), "downloading", widget.NewProgressBarInfinite(), player.GetMainWindow())
	// 	progress.Show()
	// 	DisplayErrorIfAny(scraper.AddMusicToRepository(result.VideoID(), player.GetState().Album(), result.VideoTitle()))
	// 	progress.Hide()
	// })

	// onlineBrowserDialog := dialog.NewCustom("", "( X )", container.NewBorder(
	// 	container.NewBorder(
	// 		nil,
	// 		nil,
	// 		nil,
	// 		searchButton,
	// 		searchBar,
	// 	),
	// 	nil,
	// 	nil,
	// 	nil,
	// 	scroll,
	// ), player.GetMainWindow())
	// onlineBrowserDialog.Resize(resource.GetMusicAddOnlineDialogSize())
	// return onlineBrowserDialog
}
