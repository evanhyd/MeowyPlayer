package view

import (
	"meowyplayer/model"
	"meowyplayer/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainPage struct {
	widget.BaseWidget
	homePage          *HomePage
	albumPage         *AlbumPage
	musicPage         *MusicPage
	settingPage       *SettingPage
	controllerPage    *ControllerPage
	downloadingDialog dialog.Dialog
}

func newMainPanel() *MainPage {
	p := MainPage{
		homePage:          newHomePage(),
		albumPage:         newAlbumPage(),
		musicPage:         newMusicPage(),
		settingPage:       newSettingPage(),
		controllerPage:    newControllerPage(),
		downloadingDialog: dialog.NewCustomWithoutButtons(resource.DownloadText(), widget.NewProgressBarInfinite(), getWindow()),
	}
	model.StorageClient().OnViewFocused().AttachFunc(p.showView)
	model.StorageClient().OnMusicSyncActivated().AttachFunc(p.showDownloadingDialog)

	p.ExtendBaseWidget(&p)
	return &p
}

func (p *MainPage) CreateRenderer() fyne.WidgetRenderer {
	p.musicPage.Hide()
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(resource.HomeText(), theme.HomeIcon(), p.homePage),
		container.NewTabItemWithIcon(resource.CollectionsText(), resource.CollectionTabIcon(), container.NewStack(p.albumPage, p.musicPage)),
		container.NewTabItemWithIcon(resource.SettingsText(), theme.SettingsIcon(), p.settingPage),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.SelectIndex(1)
	return widget.NewSimpleRenderer(container.NewBorder(nil, p.controllerPage, nil, nil, tabs))
}

func (p *MainPage) showView(viewID model.ViewID) {
	switch viewID {
	case model.KAlbumView:
		p.albumPage.Show()
		p.musicPage.Hide()
	case model.KMusicView:
		p.albumPage.Hide()
		p.musicPage.Show()
	}
}

func (p *MainPage) showDownloadingDialog(activated bool) {
	if activated {
		p.downloadingDialog.Show()
	} else {
		p.downloadingDialog.Hide()
	}
}
