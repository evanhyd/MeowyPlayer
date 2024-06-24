package view

import (
	"fmt"
	"playground/browser"
	"playground/resource"
	"playground/view/internal/cwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type HomePage struct {
	widget.BaseWidget
	list *cwidget.SearchList[browser.Result, *cwidget.ThumbnailCard]

	browser browser.Browser
}

func newHomePage() *HomePage {
	var p HomePage
	p = HomePage{
		list: cwidget.NewSearchList(
			container.NewVBox(),
			cwidget.NewThumbnailCardConstructor(p.showDownloadMenu, p.showDownloadMenu, p.showDownloadMenu),
			nil,
			p.searchTitle,
		),
	}

	//menu and toolbar
	p.list.AddDropDown(cwidget.NewMenuItem("YouTube", resource.YouTubeIcon, func() { p.browser = browser.NewYouTubeBrowser() }))
	p.list.AddToolbar(cwidget.NewDropDown())

	p.ExtendBaseWidget(&p)
	return &p
}

func (v *HomePage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.list)
}

func (v *HomePage) searchTitle(title string) {
	attempts := 0

	progress := widget.NewProgressBar()
	progress.TextFormatter = func() string {
		return fmt.Sprintf("%v / %v %v", attempts, resource.KSearchAttempts, resource.KAttempts)
	}
	waitDialog := dialog.NewCustomWithoutButtons(resource.KSearching, progress, getWindow())
	waitDialog.Show()
	defer waitDialog.Hide()

	for ; attempts < resource.KSearchAttempts; attempts++ {
		progress.SetValue(float64(attempts) / resource.KSearchAttempts)
		results, err := v.browser.Search(title)
		if err != nil {
			fyne.LogError("browser search failed", err)
		}
		if len(results) > 0 {
			v.list.Update(results)
			return
		}
	}
}

func (v *HomePage) showDownloadMenu(result browser.Result) {
	fmt.Println(result)
}
