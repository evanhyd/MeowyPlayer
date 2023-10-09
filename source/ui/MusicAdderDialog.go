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

func showAddOnlineMusicDialog() {
	var scraper scraper.VideoScraper = scraper.NewClipzagScraper()

	changeScraperButton := cwidget.NewButtonWithIcon("", resource.DefaultIcon(), func() { showScraperSelectionDialog(&scraper) })
	searchBar := widget.NewEntry()
	searchBar.OnSubmitted = func(title string) {
		result, err := scraper.Search(title)
		assert.NoErr(err, "failed to scrape the video info")
		fmt.Println(result)
	}
	searchButton := cwidget.NewButtonWithIcon("", theme.SearchIcon(), func() { searchBar.OnSubmitted(searchBar.Text) })

	onlineMusicDialog := dialog.NewCustom("Online Music", "( X )", container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			changeScraperButton,
			searchButton,
			searchBar,
		),
		nil,
		nil,
		nil,
	), getWindow())

	onlineMusicDialog.Resize(getWindow().Canvas().Size())
	onlineMusicDialog.Show()

	// youtubeResultData := cbinding.MakeYoutubeResultDataList()
	// scraper := scraper.NewClipzagScraper()

	// searchBar := widget.NewEntry()
	// searchBar.SetPlaceHolder(">>")
	// searchBar.OnSubmitted = func(s string) {}
	// searchButton := cwidget.NewButtonWithIcon("", resource.MusicAdderOnlineIcon(), func() { searchBar.OnSubmitted(searchBar.Text) })

	// viewList := cwidget.NewYoutubeResultView()

	// scroll := cwidget.NewList(
	// 	func() fyne.CanvasObject {
	// 		card := widget.NewCard("", "", nil)
	// 		card.Image = canvas.NewImageFromResource(defaultIcon)
	// 		card.Image.SetMinSize(resource.GetThumbnailIconSize())

	// 		videoTitle := widget.NewLabel("")
	// 		videoTitle.TextStyle = fyne.TextStyle{Bold: true, Monospace: true, Symbol: true}

	// 		videoInfo := widget.NewLabel("")
	// 		description := widget.NewLabel("")

	// 		return container.NewBorder(
	// 			nil,
	// 			nil,
	// 			card,
	// 			nil,
	// 			container.NewGridWithRows(3, videoTitle, videoInfo, description),
	// 		)
	// 	},

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

	// searchButton.OnTapped = func() {
	// 	result, err := scraper.GetSearchResult(searchBar.Text)
	// 	if err != nil {
	// 		DisplayErrorIfAny(err)
	// 		return
	// 	}
	// 	scroll.SetItems(result)
	// }

	// searchBar.OnSubmitted = func(title string) {
	// 	searchButton.OnTapped()
	// }

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

func showScraperSelectionDialog(scraper *scraper.VideoScraper) {
	scraperSelectionDialog := dialog.NewCustom("Source", "( X )", nil, getWindow())
	scraperSelectionDialog.Show()
}
