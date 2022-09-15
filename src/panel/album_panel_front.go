package panel

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/resource"
)

func NewAlbumPanel(panelInfo *custom_canvas.PanelInfo, menu *container.AppTabs, musicPanel *container.TabItem) *container.TabItem {

	//album search list
	panelInfo.AlbumSearchList = custom_canvas.NewSearchList(
		"Enter album's name...",

		SatisfyAlbumInfo,

		func(data *custom_canvas.AlbumInfo) fyne.CanvasObject {

			albumCard, err := custom_canvas.NewAlbumCard(
				data.Title, fmt.Sprintf("Music Count: %v", data.MusicCount), fyne.NewSize(ALBUM_CARD_WIDTH, ALBUM_CARD_HEIGHT),

				//add to the music list
				func() {
					panelInfo.SelectedAlbumInfo = data
					if err := LoadMusicFromAlbum(panelInfo); err != nil {
						log.Println(err)
					}
					menu.Select(musicPanel)
				},

				//secondary tap, perhaps create a menu?
				func() {
					log.Println("secondary tap")
				},
			)

			if err != nil {
				log.Println(err)
			}
			return &albumCard.Container
		},
	)

	//album adder
	albumAdderIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("album_adder_icon.png"))
	if err != nil {
		log.Println(err)
	}
	albumAdder := widget.NewButtonWithIcon("+", albumAdderIcon, func() { ShowAlbumAdderWin(panelInfo) })
	albumFrame := container.NewBorder(albumAdder, nil, nil, nil, &panelInfo.AlbumSearchList.Container)

	//Create tab
	tabIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("album_tab.png"))
	if err != nil {
		log.Println(err)
	}
	tab := container.NewTabItemWithIcon("Album", tabIcon, albumFrame)
	return tab
}

//Display adding album window
func ShowAlbumAdderWin(panelInfo *custom_canvas.PanelInfo) {

	//album adder window
	albumAdderWin := fyne.CurrentApp().NewWindow("Add Album")
	albumAdderWin.CenterOnScreen()

	//album card preview
	imagePath := resource.GetImagePath("default_album_icon.png")
	image := canvas.NewImageFromFile(imagePath)
	image.SetMinSize(fyne.NewSize(ALBUM_CARD_WIDTH, ALBUM_CARD_WIDTH))

	//album title
	title := widget.NewEntry()
	title.SetPlaceHolder("Title...")

	//image upload button
	uploadBtn := widget.NewButton("Upload Image", func() {
		ShowAddImageWin(&imagePath, image, &albumAdderWin)
	})

	//confirmation button
	confirmBtn := widget.NewButton("create!", func() {
		if err := CreateAlbumFolder(imagePath, title.Text); err != nil {
			title.SetText("")
			title.SetPlaceHolder("duplicated name")
		} else {
			LoadAlbums(panelInfo)
			albumAdderWin.Close()
		}
	})

	albumAdderWin.SetContent(
		container.NewBorder(
			container.NewVBox(image, title), nil, nil, nil,
			container.NewVBox(uploadBtn, confirmBtn),
		),
	)
	albumAdderWin.Show()
}

//Display uploading image dialog
func ShowAddImageWin(imagePath *string, image *canvas.Image, parent *fyne.Window) {

	//create image uploading window
	win := fyne.CurrentApp().NewWindow("Select Icon")

	dia := dialog.NewFileOpen(
		func(res fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println(err)
			}
			if res != nil {
				*imagePath = res.URI().Path()
				*image = *canvas.NewImageFromFile(*imagePath)
				image.SetMinSize(fyne.NewSize(ALBUM_CARD_WIDTH, ALBUM_CARD_WIDTH))
				(*parent).Content().Refresh()
			}
			win.Close()
		},
		win,
	)
	dia.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".bmp"}))

	//display
	sz := fyne.NewSize(UPLOAD_DIALOG_WIN_WIDTH, UPLOAD_DIALOG_WIN_HEIGHT)
	dia.Resize(sz)
	win.Resize(sz)
	dia.Show()
	win.Show()
}
