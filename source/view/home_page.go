package view

import (
	"fmt"
	"meowyplayer/browser"
	"meowyplayer/model"
	"meowyplayer/view/internal/cwidget"
	"meowyplayer/view/internal/resource"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	kSearchAttempts = 5
)

type HomePage struct {
	widget.BaseWidget
	searchBar *cwidget.SearchBar[[]browser.Result]
	browser   browser.Browser
}

func newHomePage() *HomePage {
	var p HomePage
	p = HomePage{
		searchBar: cwidget.NewSearchBar(
			cwidget.NewCustomList(container.NewVBox(), cwidget.NewThumbnailCardConstructor(p.onInstantPlay, p.showDownloadMenu)),
			nil,
			p.searchTitle,
		),
	}

	//menu and toolbar
	p.searchBar.AddDropDown(cwidget.NewMenuItem("YouTube", resource.YouTubeIcon(), func() { p.browser = browser.NewYouTubeBrowser() }))
	p.searchBar.AddToolbar(cwidget.NewDropDown())
	p.ExtendBaseWidget(&p)
	return &p
}

func (p *HomePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.searchBar)
}

func (p *HomePage) searchTitle(title string) {
	attempts := 0

	progress := widget.NewProgressBar()
	progress.TextFormatter = func() string {
		return fmt.Sprintf("%v / %v %v", attempts, kSearchAttempts, resource.AttemptsText())
	}
	waitDialog := dialog.NewCustomWithoutButtons(resource.SearchingText(), progress, getWindow())
	waitDialog.Show()
	defer waitDialog.Hide()

	for ; attempts < kSearchAttempts; attempts++ {
		progress.SetValue(float64(attempts) / kSearchAttempts)
		results, err := p.browser.Search(title)
		if err != nil {
			fyne.LogError("browser searchTitle failed", err)
		}
		if len(results) > 0 {
			p.searchBar.Update(results)
			return
		}
	}
}

func (p *HomePage) onInstantPlay(result browser.Result) {
	url, err := url.Parse(fmt.Sprintf("https://www.youtube.com/watch?v=%v", result.VideoID))
	if err != nil {
		fyne.LogError("instantPlay parsing URL failed", err)
	}
	err = fyne.CurrentApp().OpenURL(url)
	if err != nil {
		fyne.LogError("instantPlay open in browser failed", err)
	}
}

func (p *HomePage) showDownloadMenu(result browser.Result) {
	albums, err := model.UIClient().GetAllAlbums()
	if err != nil {
		fyne.LogError("download menu can't albums", err)
	}

	options := make([]string, 0, len(albums))
	for _, album := range albums {
		options = append(options, album.Title())
	}

	selects := widget.NewSelect(options, nil)
	selects.PlaceHolder = resource.SelectAlbumText()
	dialog.ShowCustomConfirm(resource.DownloadText(), resource.DownloadText(), resource.CancelText(), selects, func(confirm bool) {
		if index := selects.SelectedIndex(); index != -1 && confirm {
			go p.onDownload(albums[index].Key(), result)
		}
	}, getWindow())
}

func (p *HomePage) onDownload(key model.AlbumKey, result browser.Result) {
	waitDialog := dialog.NewCustomWithoutButtons(resource.DownloadText(), widget.NewProgressBarInfinite(), getWindow())
	waitDialog.Show()
	defer waitDialog.Hide()

	readCloser, err := p.browser.Download(&result)
	if err != nil {
		fyne.LogError("download failed", err)
	}
	defer readCloser.Close()
	model.UIClient().AddMusicToAlbum(key, result, readCloser)
}
