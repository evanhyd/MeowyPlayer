package ui

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/core/client"
	"meowyplayer.com/core/player"
	"meowyplayer.com/core/resource"
	"meowyplayer.com/core/ui/cbinding"
	"meowyplayer.com/core/ui/cwidget"
	"meowyplayer.com/utility/network/fileformat"
	"meowyplayer.com/utility/network/scraper"
	"meowyplayer.com/utility/pattern"
)

func newMusicTab() *container.TabItem {
	var selectedMusic resource.Music
	data := cbinding.MakeMusicDataList()
	client.Manager().AddFocusedAlbumListener(&data)

	// deleting music dialog
	deleteDialog := dialog.NewConfirm("", "Do you want to delete the music?", func(confirm bool) {
		if confirm {
			showErrorIfAny(client.Manager().DeleteMusic(client.Manager().FocusedAlbum(), selectedMusic))
		}
	}, getWindow())

	// music views
	musicViews := cwidget.NewViewList(&data, container.NewVBox(),
		func(music resource.Music) fyne.CanvasObject {
			view := cwidget.NewMusicView(&music)
			view.OnTapped = func(*fyne.PointEvent) {
				client.Manager().SetPlayList(player.MakePlayList(data.MusicList(), &music))
			}
			view.OnTappedSecondary = func(*fyne.PointEvent) {
				selectedMusic = music
				deleteDialog.Show()
			}
			return view
		},
	)

	// search bar
	searchBar := widget.NewEntry()
	searchBar.OnChanged = func(title string) {
		title = strings.ToLower(title)
		data.SetFilter(func(a resource.Music) bool {
			return strings.Contains(strings.ToLower(a.SimpleTitle()), title)
		})
	}
	client.Manager().AddFocusedAlbumListener(pattern.MakeCallback(func(resource.Album) { searchBar.SetText("") }))

	// title sorting button
	ascendTitle := -1
	titleButton := cwidget.NewButton("Title", func() {
		ascendTitle = -ascendTitle
		data.SetSorter(func(a1, a2 resource.Music) int {
			return strings.Compare(strings.ToLower(a1.Title), strings.ToLower(a2.Title)) * ascendTitle
		})
	})

	// date sorting button
	ascendDate := 1
	dateButton := cwidget.NewButton("Date", func() {
		ascendDate = -ascendDate
		data.SetSorter(func(a1, a2 resource.Music) int {
			return a1.Date.Compare(a2.Date) * ascendDate
		})
	})
	defer dateButton.OnTapped()

	// local music explorer
	localDialog := newLocalMusicDialog()

	// online music explorer
	onlineDialog := newOnlineMusicDialog()

	// add music buttons
	addLocalMusicButton := cwidget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		localDialog.Resize(getWindow().Canvas().Size())
		localDialog.Show()
	})
	addOnlineMusicButton := cwidget.NewButtonWithIcon("", resource.MusicAdderOnlineIcon, func() {
		onlineDialog.Resize(getWindow().Canvas().Size())
		onlineDialog.Show()
	})

	return container.NewTabItemWithIcon("Music", resource.MusicTabIcon, container.NewBorder(
		container.NewBorder(
			nil,
			container.NewGridWithRows(1, titleButton, dateButton),
			nil,
			container.NewGridWithRows(1, addLocalMusicButton, addOnlineMusicButton),
			searchBar,
		),
		nil,
		nil,
		nil,
		musicViews,
	))
}

func newLocalMusicDialog() dialog.Dialog {
	d := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err != nil {
			showErrorIfAny(err)
		} else if result != nil {
			showErrorIfAny(client.AddMusicFromURIReader(client.Manager().FocusedAlbum(), result))
		}
	}, getWindow())
	d.SetFilter(storage.NewExtensionFileFilter([]string{".mp3"}))
	d.SetConfirmText("Add")
	return d
}

func newOnlineMusicDialog() dialog.Dialog {
	//scraper menu
	var videoScraper scraper.VideoScraper

	//platform selector
	platformMenu := cwidget.NewDropDown("", resource.DefaultIcon)
	platformMenu.Add("YouTube", resource.YouTubeIcon, func() { videoScraper = scraper.NewClipzagScraper() })
	platformMenu.Add("BiliBili", resource.BiliBiliIcon, func() { log.Println("not implemented...") })
	platformMenu.Select(0)

	//video result view list
	videoResultData := pattern.Data[[]fileformat.VideoResult]{}
	videoResultViews := cwidget.NewViewList(&videoResultData, container.NewVBox(),
		func(video fileformat.VideoResult) fyne.CanvasObject {
			return cwidget.NewVideoView(&video, fyne.NewSize(207, 128), func() {
				showErrorIfAny(client.DownloadMusic(client.Manager().FocusedAlbum(), &video))
			})
		},
	)

	//search bar
	searchBar := widget.NewEntry()
	searchBar.SetPlaceHolder("Search Video")
	searchBar.ActionItem = cwidget.NewButtonWithIcon("", theme.SearchIcon(), func() { searchBar.OnSubmitted(searchBar.Text) })
	searchBar.OnSubmitted = func(title string) {
		videos, err := videoScraper.Search(title)
		if err != nil {
			showErrorIfAny(err)
			return
		}
		videoResultData.Set(videos)
	}

	d := dialog.NewCustom("", "X", container.NewBorder(
		container.NewBorder(nil, nil, platformMenu, nil, searchBar),
		nil,
		nil,
		nil,
		videoResultViews,
	), getWindow())
	return d
}
