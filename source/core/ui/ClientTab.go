package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/pattern"
)

func newCollectionViewList(data pattern.Subject[[]resource.CollectionInfo], onDownload func(*resource.CollectionInfo)) *cwidget.ViewList[resource.CollectionInfo] {
	return cwidget.NewViewList(data, container.NewVBox(),
		func(info resource.CollectionInfo) fyne.CanvasObject {
			return cwidget.NewCollectionInfoView(&info, onDownload)
		},
	)
}

func newClientTab() *container.TabItem {
	account := resource.Account{Name: "UnboxTheCat", ID: "0x00000000"}
	accountView := cwidget.NewAccountView()
	accountView.SetAccount(&account)
	accountView.SetOnAutoBackup(func(b bool) { log.Println("backup:", b) })

	//request download from server
	requestDownload := func(info *resource.CollectionInfo) {
		showErrorIfAny(client.Manager().ClientRequestDownload(&account, info))
	}

	infoData := pattern.Data[[]resource.CollectionInfo]{}
	infoViewList := newCollectionViewList(&infoData, requestDownload)

	//request list from server
	requestList := func(url string) {
		client.Config().SetServer(url)
		if infos, err := client.Manager().ClientRequestList(&account); err != nil {
			showErrorIfAny(err)
		} else {
			infoData.Set(infos)
		}
	}

	//request upload to server
	requestUpload := func() {
		showErrorIfAny(client.Manager().ClientRequestUpload(&account))
	}

	serverEntry := widget.NewEntry()
	serverEntry.ActionItem = cwidget.NewButtonWithIcon("", theme.ComputerIcon(), func() { serverEntry.OnSubmitted(serverEntry.Text) })
	serverEntry.OnSubmitted = requestList
	serverEntry.SetPlaceHolder("server url")
	serverEntry.SetText(client.Config().ServerUrl)

	uploadButton := cwidget.NewButtonWithIcon("upload", theme.UploadIcon(), requestUpload)

	return container.NewTabItemWithIcon("Client", theme.AccountIcon(), container.NewGridWithColumns(2,
		accountView,
		container.NewBorder(serverEntry, uploadButton, nil, nil, infoViewList),
	))
}
