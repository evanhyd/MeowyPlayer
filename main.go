package main

import (
	"bufio"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/panel"
	"meowyplayer.com/src/resource"
	"meowyplayer.com/src/seeker"
)

func main() {

	//create necessary directories
	os.Mkdir(resource.GetMusicBasePath(), os.ModePerm)
	os.Mkdir(resource.GetAlbumBasePath(), os.ModePerm)

	//Loading main window
	log.Println("loading main window...")
	fyne.SetCurrentApp(app.New())
	mainWin := fyne.CurrentApp().NewWindow("Meowy Player")
	mainWin.Resize(fyne.NewSize(340.0, 340.0*1.618))
	mainWin.CenterOnScreen()
	icon, err := fyne.LoadResourceFromPath(resource.GetImagePath("main_window_icon.png"))
	if err != nil {
		log.Panic(err)
	}
	mainWin.SetIcon(icon)

	//load components
	log.Println("loading menu...")
	menu := container.NewAppTabs()
	menu.SetTabLocation(container.TabLocationLeading)
	panelInfo := custom_canvas.NewPanelInfo()
	musicPanel := panel.NewMusicPanel(panelInfo)
	albumPanel := panel.NewAlbumPanel(panelInfo, menu, musicPanel)
	menu.Append(albumPanel)
	menu.Append(musicPanel)

	//load album on launch
	err = panel.LoadAlbums(panelInfo)
	if err != nil {
		log.Println(err)
	}

	//load seeker
	log.Println("loading seeker...")
	seekerUI := seeker.NewSeekerUI()

	//combine
	mainWin.SetContent(container.NewBorder(nil, seekerUI, nil, nil, menu))
	mainWin.ShowAndRun()

	log.Println("checking unused file...")
	if err := RemoveUnusedMusic(); err != nil {
		log.Println(err)
	}
}

func RemoveUnusedMusic() error {

	//read album directories
	albumDirs, err := os.ReadDir(resource.GetAlbumBasePath())
	if err != nil {
		return err
	}

	//read music config
	inUse := make(map[string]struct{})
	for _, albumDir := range albumDirs {

		if albumDir.IsDir() {
			configFile, err := os.Open(resource.GetAlbumConfigPath(albumDir.Name()))
			if err != nil {
				return err
			}
			defer configFile.Close()

			//read music name from config file
			scanner := bufio.NewScanner(configFile)
			for line := 0; scanner.Scan(); line++ {
				inUse[scanner.Text()] = struct{}{}
			}
		}
	}

	//remove unused music
	musicDirs, err := os.ReadDir(resource.GetMusicBasePath())
	if err != nil {
		return err
	}
	for _, musicDir := range musicDirs {
		if !musicDir.IsDir() {
			if _, ok := inUse[musicDir.Name()]; !ok {
				if err := os.Remove(resource.GetMusicPath(musicDir.Name())); err != nil {
					return err
				}
				log.Printf("removed %v\n", musicDir.Name())
			}
		}
	}

	return nil
}
