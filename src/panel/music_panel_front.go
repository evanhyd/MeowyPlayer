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
	var musicAdderWin fyne.Window = nil
	musicAdder := widget.NewButtonWithIcon("+", musicAdderIcon, func() {
		if panelInfo.SelectedAlbumInfo != nil {
			if musicAdderWin == nil {
				musicAdderWin = GetMusicAdderWin(panelInfo)
				musicAdderWin.SetOnClosed(func() { musicAdderWin = nil })
				musicAdderWin.Show()
			} else {
				musicAdderWin.RequestFocus()
			}
		}
	})

	//music search list
	var musicRemoverWin fyne.Window = nil
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
				func() {
					if musicRemoverWin == nil {
						musicRemoverWin = GetMusicRemoverWin(panelInfo, data.Index)
						musicRemoverWin.SetOnClosed(func() { musicRemoverWin = nil })
						musicRemoverWin.Show()
					}
				},
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

//Show music removing window
func GetMusicRemoverWin(panelInfo *custom_canvas.PanelInfo, musicIndex int) fyne.Window {

	musicInfo := panelInfo.MusicSearchList.DataList[musicIndex]

	win := fyne.CurrentApp().NewWindow("Remove Music")
	win.CenterOnScreen()

	dia := container.NewBorder(
		widget.NewLabel("Remove \""+musicInfo.Title+"\" ?"),
		nil,
		nil,
		nil,

		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Yes", func() {
				if err := removeMusic(panelInfo.SelectedAlbumInfo.Title, musicIndex); err != nil {
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
	return win
}

//Display adding music window
func GetMusicAdderWin(panelInfo *custom_canvas.PanelInfo) fyne.Window {
	addLocalMusicBtn := widget.NewButton("From Local", func() { ShowAddLocalMusicWin(panelInfo) })
	addURLMusicBtn := widget.NewButton("From URL", func() { ShowAddURLMusicWin(panelInfo) })
	addYoutubeMusicBtn := widget.NewButton("From Youtube", func() { ShowAddYoutubeMusicWin(panelInfo) })
	win := fyne.CurrentApp().NewWindow("Add Music")
	win.SetContent(container.NewVBox(addLocalMusicBtn, addURLMusicBtn, addYoutubeMusicBtn))
	win.CenterOnScreen()
	return win
}

//Create add local music window
func ShowAddLocalMusicWin(panelInfo *custom_canvas.PanelInfo) {

	win := fyne.CurrentApp().NewWindow("From Local")
	win.CenterOnScreen()

	dia := dialog.NewFileOpen(
		func(res fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println(err)
			}
			if res != nil {

				//add music from local
				musicTitle := res.URI().Path()[strings.LastIndex(res.URI().Path(), `/`)+1:]
				if err := addMusicFromLocal(res.URI().Path(), musicTitle, panelInfo.SelectedAlbumInfo.Title); err != nil {
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

//Create add url music window
func ShowAddURLMusicWin(panelInfo *custom_canvas.PanelInfo) {

	win := fyne.CurrentApp().NewWindow("From URL")
	win.CenterOnScreen()

	titleEntry := widget.NewEntry()
	urlEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title...")
	urlEntry.SetPlaceHolder("Youtube URL...")

	downloadBtn := widget.NewButton("Download", func() {

		//download music
		if err := addMusicFromURL(urlEntry.Text, titleEntry.Text+`.mp3`, panelInfo.SelectedAlbumInfo.Title); err != nil {
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

//Create add youtube music window
func ShowAddYoutubeMusicWin(panelInfo *custom_canvas.PanelInfo) {

	win := fyne.CurrentApp().NewWindow("From Youtube")
	win.Resize(fyne.NewSize(340.0, 340.0*1.618))
	win.CenterOnScreen()

	queryEntry := widget.NewEntry()
	searchBtn := widget.NewButton("Search", func() {})

	win.SetContent(
		container.NewBorder(
			container.NewBorder(nil, nil, nil, searchBtn, queryEntry),
			nil,
			nil,
			nil,
		),
	)
	win.Show()
}

/*

<div class=\"thumbnail\"> <a class=\"item-thumbnail\" style=\"position: relative\"

{
    "status": "success",
    "result": "<div class=\"row\" id=\"search-result\"> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/FvSWwYk7G3k\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"FUNNY CAT MEMES COMPILATION OF 2022 PART 57\" src=\"https:\/\/i.ytimg.com\/vi\/FvSWwYk7G3k\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/FvSWwYk7G3k\" target=\"_blank\"> FUNNY CAT MEMES COMPILATION OF 2022 PART 57 <\/a> <p> <a href=\"\/youtube\/FvSWwYk7G3k\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/katD5xvV2t8\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Cat crying\" src=\"https:\/\/i.ytimg.com\/vi\/katD5xvV2t8\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/katD5xvV2t8\" target=\"_blank\"> Cat crying <\/a> <p> <a href=\"\/youtube\/katD5xvV2t8\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/WyN-E60PJh0\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Funny Animal Videos 2022 \ud83d\ude02 - Best Dogs And Cats Videos \ud83d\ude3a\ud83d\ude0d #23\" src=\"https:\/\/i.ytimg.com\/vi\/WyN-E60PJh0\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/WyN-E60PJh0\" target=\"_blank\"> Funny Animal Videos 2022 \ud83d\ude02 - Best Dogs And Cats Videos \ud83d\ude3a\ud83d\ude0d #23 <\/a> <p> <a href=\"\/youtube\/WyN-E60PJh0\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/zFXMOfFzoRc\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Funniest Cat Videos on the Planet #3 - Funny Cats and Dogs Videos\" src=\"https:\/\/i.ytimg.com\/vi\/zFXMOfFzoRc\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/zFXMOfFzoRc\" target=\"_blank\"> Funniest Cat Videos on the Planet #3 - Funny Cats and Dogs Videos <\/a> <p> <a href=\"\/youtube\/zFXMOfFzoRc\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/H24Epu-1k6c\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Cat TV for Cats to Watch \ud83d\ude3a Halloween pumpkin, feral squirrels, blackbirds \ud83d\udc3f 8 Hours(4K HDR)\" src=\"https:\/\/i.ytimg.com\/vi\/H24Epu-1k6c\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/H24Epu-1k6c\" target=\"_blank\"> Cat TV for Cats to Watch \ud83d\ude3a Halloween pumpkin, feral squirrels, blackbirds \ud83d\udc3f 8 Hours(4K HDR) <\/a> <p> <a href=\"\/youtube\/H24Epu-1k6c\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/ODiHIS3AWLg\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"BEST CAT TIKTOKS!! #4\" src=\"https:\/\/i.ytimg.com\/vi\/ODiHIS3AWLg\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/ODiHIS3AWLg\" target=\"_blank\"> BEST CAT TIKTOKS!! #4 <\/a> <p> <a href=\"\/youtube\/ODiHIS3AWLg\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/xbs7FT7dXYc\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Videos for Cats to Watch - 8 Hour Birds Bonanza - Cat TV Bird Watch\" src=\"https:\/\/i.ytimg.com\/vi\/xbs7FT7dXYc\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/xbs7FT7dXYc\" target=\"_blank\"> Videos for Cats to Watch - 8 Hour Birds Bonanza - Cat TV Bird Watch <\/a> <p> <a href=\"\/youtube\/xbs7FT7dXYc\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/Pms2R85EWC4\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Kitten Kiki feels annoyed when mother cat is around, but feels lonely when she isn't...\" src=\"https:\/\/i.ytimg.com\/vi\/Pms2R85EWC4\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/Pms2R85EWC4\" target=\"_blank\"> Kitten Kiki feels annoyed when mother cat is around, but feels lonely when she isn't... <\/a> <p> <a href=\"\/youtube\/Pms2R85EWC4\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/fXuYx81aYoI\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Cat Goes On Walks With His Dog In The Cutest Way | The Dodo Odd Couples\" src=\"https:\/\/i.ytimg.com\/vi\/fXuYx81aYoI\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/fXuYx81aYoI\" target=\"_blank\"> Cat Goes On Walks With His Dog In The Cutest Way | The Dodo Odd Couples <\/a> <p> <a href=\"\/youtube\/fXuYx81aYoI\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/nbFTp0RadzQ\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"I Taught My Cat to Play Minecraft\" src=\"https:\/\/i.ytimg.com\/vi\/nbFTp0RadzQ\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/nbFTp0RadzQ\" target=\"_blank\"> I Taught My Cat to Play Minecraft <\/a> <p> <a href=\"\/youtube\/nbFTp0RadzQ\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/Y0NWybTQv9A\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Funny animals - Funny cats \/ dogs - Funny animal videos 232\" src=\"https:\/\/i.ytimg.com\/vi\/Y0NWybTQv9A\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/Y0NWybTQv9A\" target=\"_blank\"> Funny animals - Funny cats \/ dogs - Funny animal videos 232 <\/a> <p> <a href=\"\/youtube\/Y0NWybTQv9A\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <div class=\"col-xs-6 col-sm-4 col-md-3\">
		href=\"\/youtube\/S5JqSlAsldQ\" target=\"_blank\"> <img class=\"ythumbnail\" alt=\"Real Meanings Behind 9 Strange Cat Behaviors Explained\" src=\"https:\/\/i.ytimg.com\/vi\/S5JqSlAsldQ\/0.jpg\"> <\/a> <div class=\"search-info\"> <a href=\"\/youtube\/S5JqSlAsldQ\" target=\"_blank\"> Real Meanings Behind 9 Strange Cat Behaviors Explained <\/a> <p> <a href=\"\/youtube\/S5JqSlAsldQ\" target=\"_blank\" class=\"btn btn-success btn-xs\"> <i class=\"glyphicon glyphicon-download-alt\"><\/i>&nbsp; Download video<\/a> <\/p> <br \/> <\/div> <\/div> <\/div> <\/div> <script type=\"text\/javascript\"> k_data_vid = \"cat\"; video_service = \"youtube\"; video_extractor = \"search\"; <\/script> "
}


*/
