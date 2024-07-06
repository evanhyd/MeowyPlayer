package view

import (
	"playground/model"
	"playground/view/internal/cwidget"
	"playground/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainPanel struct {
	widget.BaseWidget
	homePage    *HomePage
	albumPage   *AlbumPage
	musicPage   *MusicPage
	settingPage *SettingPage
	controller  *cwidget.MediaController
}

func newMainPanel() *MainPanel {
	p := MainPanel{
		homePage:    newHomePage(),
		albumPage:   newAlbumPage(),
		musicPage:   newMusicPage(),
		settingPage: newSettingPage(),
	}
	model.Instance().OnAlbumViewFocused().AttachFunc(p.showAlbumTab)
	model.Instance().OnMusicViewFocused().AttachFunc(p.showMusicTab)

	//TODO: create music controller
	p.ExtendBaseWidget(&p)
	return &p
}

func (p *MainPanel) CreateRenderer() fyne.WidgetRenderer {
	p.musicPage.Hide()
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(resource.KHomeText, theme.HomeIcon(), p.homePage),
		container.NewTabItemWithIcon(resource.KCollectionText, resource.CollectionTabIcon, container.NewStack(p.albumPage, p.musicPage)),
		container.NewTabItemWithIcon(resource.KSettingText, theme.SettingsIcon(), p.settingPage),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	return widget.NewSimpleRenderer(container.NewBorder(nil, widget.NewLabel("Controller"), nil, nil, tabs))
}

func (p *MainPanel) showAlbumTab(bool) {
	p.albumPage.Show()
	p.musicPage.Hide()
}

func (p *MainPanel) showMusicTab(bool) {
	p.albumPage.Hide()
	p.musicPage.Show()
}
