package view

import (
	"playground/model"
	"playground/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainPanel struct {
	widget.BaseWidget
	tabs *container.AppTabs
}

func getWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func NewMainPanel(client *model.Client) *MainPanel {
	v := &MainPanel{}

	//create tabs
	v.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon(resource.KHomeText, theme.HomeIcon(), NewHomePage(client)),                 //home
		container.NewTabItemWithIcon(resource.KCollectionText, resource.CollectionTabIcon, NewSwitcher(client)), //collection
		container.NewTabItemWithIcon(resource.KAccountText, theme.AccountIcon(), NewAccountPage(client)),        //account
		container.NewTabItemWithIcon(resource.KSettingText, theme.SettingsIcon(), NewSettingPage(client)),       //setting
	)
	v.tabs.SetTabLocation(container.TabLocationLeading)

	//TODO: create music controller
	v.ExtendBaseWidget(v)
	return v
}

func (v *MainPanel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, widget.NewLabel("Controller"), nil, nil, v.tabs))
}
