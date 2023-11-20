package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	userNameEntry.ActionItem = cwidget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() { userNameEntry.OnSubmitted(userNameEntry.Text) })
	userNameEntry.OnSubmitted = func(name string) { client.Config().SetName(name) }
	client.Config().AddListener(pattern.MakeCallback(func(config resource.Config) { userNameEntry.SetText(config.Name) }))

	//collection info list
	infoData := pattern.Data[[]resource.CollectionInfo]{}
	infoViewList := cwidget.NewViewList(&infoData, container.NewVBox(),
		func(info resource.CollectionInfo) fyne.CanvasObject {
			return cwidget.NewCollectionInfoView(&info, func() {
				progress := dialog.NewCustomWithoutButtons("downloading", widget.NewProgressBarInfinite(), getWindow())
				progress.Show()
				defer progress.Hide()
				showErrorIfAny(client.RequestDownload(&info))
			})
		},
	)

	//server ip
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("server url")
	serverEntry.ActionItem = cwidget.NewButtonWithIcon("", theme.ComputerIcon(), func() { serverEntry.OnSubmitted(serverEntry.Text) })
	serverEntry.OnSubmitted = func(url string) {
		progress := dialog.NewCustomWithoutButtons("listing", widget.NewProgressBarInfinite(), getWindow())
		progress.Show()
		defer progress.Hide()
		client.Config().SetServerUrl(url)
		infos, err := client.RequestList()
		if err != nil {
			showErrorIfAny(err)
			return
		}
		infoData.Set(infos)
	}
	client.Config().AddListener(pattern.MakeCallback(func(config resource.Config) { serverEntry.SetText(config.ServerUrl) }))

	//upload config
	uploadButton := cwidget.NewButtonWithIcon("upload", theme.UploadIcon(), func() {
		progress := dialog.NewCustomWithoutButtons("uploading", widget.NewProgressBarInfinite(), getWindow())
		progress.Show()
		showErrorIfAny(client.RequestUpload())
		progress.Hide()
	})

	return container.NewTabItemWithIcon("Client", theme.AccountIcon(), container.NewGridWithColumns(2,
		container.NewVBox(userNameEntry),
		container.NewBorder(serverEntry, uploadButton, nil, nil, infoViewList),
	))
}
