package panel

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/resource"
	"meowyplayer.com/src/seeker"
)

func NewMusicPanel(panelInfo *custom_canvas.PanelInfo) *container.TabItem {

	// music adder
	musicAdderIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("music_adder_icon.png"))
	if err != nil {
		log.Println(err)
	}
	musicAdder := widget.NewButtonWithIcon("+", musicAdderIcon, func() { ShowMusicAdderWin(panelInfo) })

	//music search list
	panelInfo.MusicSearchList = custom_canvas.NewSearchList(
		"Enter music's name...",

		SatisfyMusicInfo,

		func(data *custom_canvas.MusicInfo) fyne.CanvasObject {
			card, err := custom_canvas.NewMusicCard(

				//music title and count
				data.Title, data.Duration,

				//play music
				func() { seeker.TheUniquePlayer.SetPlaylist(panelInfo.MusicSearchList.DataList, data.Index) },

				//remove music
				func() { ShowMusicRemoverWin(panelInfo, data.Index) },
			)
			if err != nil {
				log.Println(err)
			}
			return &card.Container
		},
	)

	//Create tab
	tabIcon, err := fyne.LoadResourceFromPath(resource.GetImagePath("music_tab.png"))
	if err != nil {
		log.Println(err)
	}
	return container.NewTabItemWithIcon("Music", tabIcon, container.NewBorder(musicAdder, nil, nil, nil, &panelInfo.MusicSearchList.Container))
}

//Display adding music window
func ShowMusicAdderWin(panelInfo *custom_canvas.PanelInfo) {

	if panelInfo.SelectedAlbumInfo == nil {
		return
	}

	addLocalMusicBtn := widget.NewButton("From Local", func() {
		ShowAddLocalMusicWin(panelInfo)
	})

	addRemoteMusicBtn := widget.NewButton("From Remote", func() {
		ShowAddRemoteMusicWin(panelInfo)
	})

	win := fyne.CurrentApp().NewWindow("Add Music")
	win.SetContent(container.NewVBox(addLocalMusicBtn, addRemoteMusicBtn))
	win.Show()
}

//Selecting music from local window
func ShowAddLocalMusicWin(panelInfo *custom_canvas.PanelInfo) {

	//create image uploading window
	win := fyne.CurrentApp().NewWindow("From Local")
	dia := dialog.NewFileOpen(
		func(res fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println(err)
			}
			if res != nil {

				//add music from local
				musicTitle := res.URI().Path()[strings.LastIndex(res.URI().Path(), `/`)+1:]
				if err := AddMusicFromLocal(res.URI().Path(), musicTitle, panelInfo.SelectedAlbumInfo.Title); err != nil {
					log.Println(err)
				}

				//reload album list
				if err := LoadAlbumFromDir(panelInfo); err != nil {
					log.Println(err)
				}

				//reload music list
				if err := LoadMusicFromAlbum(panelInfo); err != nil {
					log.Println(err)
				}
			}
			win.Close()
		},
		win,
	)
	dia.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	dia.SetConfirmText("Upload")

	//display
	sz := fyne.NewSize(UPLOAD_DIALOG_WIN_WIDTH, UPLOAD_DIALOG_WIN_HEIGHT)
	dia.Resize(sz)
	win.Resize(sz)
	dia.Show()
	win.Show()
}

//Downloading music from remote window
func ShowAddRemoteMusicWin(panelInfo *custom_canvas.PanelInfo) {

	win := fyne.CurrentApp().NewWindow("Download Music")

	titleEntry := widget.NewEntry()
	urlEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title...")
	urlEntry.SetPlaceHolder("Youtube URL...")

	downloadBtn := widget.NewButton("Download", func() {

		//download music
		if err := AddMusicFromRemote(urlEntry.Text, titleEntry.Text+`.mp3`, panelInfo.SelectedAlbumInfo.Title); err != nil {
			log.Println(err)
		}

		//reload album list
		if err := LoadAlbumFromDir(panelInfo); err != nil {
			log.Println(err)
		}

		//reload music list
		if err := LoadMusicFromAlbum(panelInfo); err != nil {
			log.Println(err)
		}

		win.Close()
	})

	win.SetContent(container.NewBorder(container.NewVBox(titleEntry, urlEntry), nil, nil, nil, downloadBtn))
	win.Show()
}

//Show music removing window
func ShowMusicRemoverWin(panelInfo *custom_canvas.PanelInfo, musicIndex int) {

	musicInfo := panelInfo.MusicSearchList.DataList[musicIndex]

	win := fyne.CurrentApp().NewWindow("Remove Music")
	dia := container.NewBorder(
		widget.NewLabel("Remove \""+musicInfo.Title+"\" ?"),
		nil,
		nil,
		nil,

		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Yes", func() {
				if err := RemoveMusic(panelInfo.SelectedAlbumInfo.Title, musicIndex); err != nil {
					log.Println(err)
				}
				if err := LoadMusicFromAlbum(panelInfo); err != nil {
					log.Println(err)
				}
				win.Close()
			}),
			layout.NewSpacer(),
			widget.NewButton("No", func() { win.Close() }),
			layout.NewSpacer(),
		),
	)

	win.SetContent(dia)
	dia.Show()
	win.Show()
}
