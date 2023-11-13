package ui

import (
	"log"

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
	account := resource.Account{Name: "UnboxTheCat", ID: "0x00000000"}
	accountView := cwidget.NewAccountView()
	accountView.SetAccount(&account)
	accountView.SetOnAutoBackup(func(b bool) { log.Println("backup:", b) })

	//collection info list
	infoData := pattern.Data[[]resource.CollectionInfo]{}
	infoViewList := cwidget.NewViewList(&infoData, container.NewVBox(),
		func(info resource.CollectionInfo) fyne.CanvasObject {
			return cwidget.NewCollectionInfoView(&info, func(info *resource.CollectionInfo) {
				showErrorIfAny(client.Manager().ClientRequestDownload(&account, info))
			})
		},
	)

	serverEntry := widget.NewEntry()
	serverEntry.ActionItem = cwidget.NewButtonWithIcon("", theme.ComputerIcon(), func() { serverEntry.OnSubmitted(serverEntry.Text) })
	serverEntry.OnSubmitted = func(url string) {
		progress := dialog.NewCustomWithoutButtons("listing", widget.NewProgressBarInfinite(), getWindow())
		progress.Show()
		defer progress.Hide()
		client.Config().SetServer(url)
		infos, err := client.Manager().ClientRequestList(&account)
		if err != nil {
			showErrorIfAny(err)
			return
		}
		infoData.Set(infos)
	}
	serverEntry.SetPlaceHolder("server url")
	serverEntry.SetText(client.Config().ServerUrl)

	uploadButton := cwidget.NewButtonWithIcon("upload", theme.UploadIcon(), func() {
		progress := dialog.NewCustomWithoutButtons("uploading", widget.NewProgressBarInfinite(), getWindow())
		progress.Show()
		showErrorIfAny(client.Manager().ClientRequestUpload(&account))
		progress.Hide()
	})

	return container.NewTabItemWithIcon("Client", theme.AccountIcon(), container.NewGridWithColumns(2,
		accountView,
		container.NewBorder(serverEntry, uploadButton, nil, nil, infoViewList),
	))
}
