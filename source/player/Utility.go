package player

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/resource"
)

func GetMainWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

// Read album name, music config, cover icon from the given directory
func ReadAlbumsFromDisk() ([]Album, error) {
	directories, err := os.ReadDir(resource.GetAlbumRootPath())
	if err != nil {
		return nil, err
	}

	albums := []Album{}
	for _, directory := range directories {
		if directory.IsDir() {

			//read album config
			configPath := resource.GetAlbumConfigPath(directory.Name())
			config, err := os.ReadFile(configPath)
			if err != nil {
				return nil, err
			}

			//read last modified date
			info, err := os.Stat(configPath)
			if err != nil {
				return nil, err
			}

			musicNumber := bytes.Count(config, []byte{'\n'})

			coverIcon, err := fyne.LoadResourceFromPath(resource.GetAlbumIconPath(directory.Name()))
			if err != nil {
				return nil, err
			}
			albumCover := canvas.NewImageFromResource(coverIcon)
			albumCover.SetMinSize(resource.GetAlbumCoverSize())
			albums = append(albums, Album{directory.Name(), info.ModTime(), musicNumber, albumCover})
		}
	}
	return albums, nil
}

// Read a list of music from the album config
func ReadMusicFromDisk(album Album) ([]Music, error) {

	if album.IsPlaceHolder() {
		return nil, nil
	}

	// a rough estimation of the music duration in nanoseconds
	estimateDuration := func(musicFileInfo fs.FileInfo) time.Duration {
		return time.Duration(musicFileInfo.Size() * MAGIC_RATIO / (AUDIO_BIT_DEPTH * NUM_OF_CHANNELS * SAMPLING_RATE))
	}

	config, err := os.ReadFile(resource.GetAlbumConfigPath(album.title))
	if err != nil {
		return nil, err
	}

	//read music name from config
	music := []Music{}
	scanner := bufio.NewScanner(bytes.NewReader(config))
	for scanner.Scan() {
		info, err := os.Stat(resource.GetMusicPath(scanner.Text()))
		if err != nil {
			return nil, err
		}

		music = append(music, Music{scanner.Text(), estimateDuration(info), info.ModTime()})
	}

	return music, scanner.Err()
}

// Refresh album tab GUI content
func RefreshAlbumTab() error {
	albums, err := ReadAlbumsFromDisk()
	if err != nil {
		return err
	}
	state.onUpdateAlbums.NotifyAll(albums)
	return nil
}

// Refresh music tab GUI content
func RefreshMusicTab() error {
	musics, err := ReadMusicFromDisk(state.album)
	if err != nil {
		return err
	}
	state.musics = musics
	state.onUpdateMusics.NotifyAll(musics)
	return nil
}

// Refresh seeker GUI content
func RefreshSeeker(music Music) error {
	panic("not implemented yet")
}

// Switch to album tab
func FocusAlbumTab() {
	state.onFocusAlbumTab.NotifyAll()
}

// Switch to music tab
func FocusMusicTab() {
	state.onFocusMusicTab.NotifyAll()
}

// When user click the album card
func UserSelectAlbum(selectedAlbum Album) error {
	if state.album != selectedAlbum {
		state.album = selectedAlbum
		if err := RefreshMusicTab(); err != nil {
			return err
		}
	}
	FocusMusicTab()
	return nil
}

// When user click the music row
func UserSelectMusic(selectedMusic Music) {
	state.onUpdateSeeker.NotifyAll(state.album, state.musics, selectedMusic)
}

// Add new empty album with random name and icon
func AddNewAlbum() error {
	title := strconv.FormatInt(rand.Int63(), 10)

	if err := os.Mkdir(resource.GetAlbumFolderPath(title), fs.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(resource.GetAlbumConfigPath(title), []byte{}, os.ModePerm); err != nil {
		return err
	}

	iconFile, err := os.Create(resource.GetAlbumIconPath(title))
	if err != nil {
		return err
	}
	defer iconFile.Close()

	iconColor := color.NRGBA{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())}
	iconImage := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	iconImage.SetNRGBA(0, 0, iconColor)
	if err := png.Encode(iconFile, iconImage); err != nil {
		return err
	}

	return RefreshAlbumTab()
}

// Rename album's title
func RenameAlbum(album Album, newTitle string) error {
	oldPath := resource.GetAlbumFolderPath(album.Title())
	newPath := resource.GetAlbumFolderPath(newTitle)

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return fmt.Errorf("album \"%v\" is already existed", newTitle)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if state.album.Title() == album.Title() {
		state.album.title = newTitle
	}
	return RefreshAlbumTab()
}

// Set album cover icon
func SetAlbumCover(album Album, coverIconPath string) error {
	coverIconData, err := os.ReadFile(coverIconPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(resource.GetAlbumIconPath(album.Title()), coverIconData, os.ModePerm); err != nil {
		return err
	}

	return RefreshAlbumTab()
}

// Remove album
func RemoveAlbum(album Album) error {
	if err := os.RemoveAll(resource.GetAlbumFolderPath(album.Title())); err != nil {
		return err
	}

	if err := RefreshAlbumTab(); err != nil {
		return err
	}

	if state.album.Title() == album.Title() {
		state.album = GetPlaceHolderAlbum()
		return RefreshMusicTab()
	}
	return nil
}

// Add music to the album
func AddMusicToAlbum(album Album, sourcePath, musicTitle string) error {

	//add to music repository
	music, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(resource.GetMusicPath(musicTitle), music, os.ModePerm); err != nil {
		return err
	}

	//add music to the config
	config, err := os.OpenFile(resource.GetAlbumConfigPath(album.Title()), os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer config.Close()

	//check for duplicated name
	scanner := bufio.NewScanner(config)
	for scanner.Scan() {
		if scanner.Text() == musicTitle {
			return fmt.Errorf("%v is already in the album", musicTitle)
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	if _, err := config.WriteString(musicTitle + "\n"); err != nil {
		return err
	}

	//update GUI
	if err := RefreshAlbumTab(); err != nil {
		return err
	}
	return RefreshMusicTab()
}

// Remove music from the album
func RemoveMusicFromAlbum(album Album, music Music) error {
	//load music config
	albumPath := resource.GetAlbumConfigPath(album.Title())
	config, err := os.ReadFile(albumPath)
	if err != nil {
		return err
	}

	//remove music name from config
	titles := bytes.Split(config, []byte{'\n'})
	musicIndex := slices.IndexFunc(titles, func(title []byte) bool { return slices.Equal(title, []byte(music.Title())) })
	if musicIndex != -1 {
		lastMusicIndex := len(titles) - 2
		newLineIndex := len(titles) - 1
		titles[musicIndex], titles[lastMusicIndex] = titles[lastMusicIndex], titles[newLineIndex]
		titles = titles[:newLineIndex]
	}

	//override the config
	if err := os.WriteFile(albumPath, bytes.Join(titles, []byte{'\n'}), os.ModePerm); err != nil {
		return err
	}

	if err := RefreshAlbumTab(); err != nil {
		return err
	}
	return RefreshMusicTab()
}
