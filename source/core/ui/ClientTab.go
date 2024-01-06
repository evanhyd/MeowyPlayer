package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newClientTab() *container.TabItem {
	//username
	userNameEntry := widget.NewEntry()
	userNameEntry.SetPlaceHolder("username")
	client.Config().AddListener(pattern.MakeCallback(func(config resource.Config) { userNameEntry.SetText(config.Name) }))

	//password
	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("password")
	passwordEntry.Password = true

	//server ip
	var infoData pattern.Data[[]resource.CollectionInfo]
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("server url")
	serverEntry.ActionItem = cwidget.NewButtonWithIcon("", theme.ComputerIcon(), func() { serverEntry.OnSubmitted(serverEntry.Text) })
	serverEntry.OnSubmitted = func(url string) {
		loadingCall(func() {
			client.Config().SetName(userNameEntry.Text)
			client.Config().SetServerUrl(url)
			infos, err := client.RequestList(serverEntry.Text, userNameEntry.Text, passwordEntry.Text)
			if err != nil {
				showErrorIfAny(err)
				return
			}
			infoData.Set(infos)
		})()
	}

	client.Config().AddListener(pattern.MakeCallback(func(config resource.Config) { serverEntry.SetText(config.ServerUrl) }))

	//collection info list
	infoData = pattern.Data[[]resource.CollectionInfo]{}
	infoViewList := cwidget.NewViewList(&infoData, container.NewVBox(),
		func(info resource.CollectionInfo) fyne.CanvasObject {
			return cwidget.NewCollectionInfoView(&info, loadingCall(func() {
				showErrorIfAny(client.RequestDownload(serverEntry.Text, userNameEntry.Text, passwordEntry.Text, &info))
			}))
		},
	)

	//upload config
	uploadButton := cwidget.NewButtonWithIcon("upload", theme.UploadIcon(), loadingCall(func() {
		showErrorIfAny(client.RequestUpload(serverEntry.Text, userNameEntry.Text, passwordEntry.Text))
	}))

	return container.NewTabItemWithIcon("Client", theme.AccountIcon(), container.NewGridWithColumns(2,
		container.NewVBox(userNameEntry, passwordEntry),
		container.NewBorder(serverEntry, uploadButton, nil, nil, infoViewList),
	))
}
