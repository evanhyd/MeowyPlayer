package player

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/resource"
)

var state *State

func init() {
	state = NewState()
}

func GetState() *State {
	return state
}

type State struct {
	album                       Album
	musics                      []Music
	onReadAlbumsFromDiskSubject pattern.OneArgSubject[[]Album]
	onSelectAlbumSubject        pattern.OneArgSubject[Album]
	onReadMusicsDiskSubject     pattern.OneArgSubject[[]Music]
	onSelectMusicSubject        pattern.ThreeArgSubject[Album, []Music, Music]
}

func NewState() *State {
	return &State{}
}

func (state *State) OnReadAlbumsFromDiskSubject() *pattern.OneArgSubject[[]Album] {
	return &state.onReadAlbumsFromDiskSubject
}

func (state *State) OnSelectAlbumSubject() *pattern.OneArgSubject[Album] {
	return &state.onSelectAlbumSubject
}

func (state *State) OnReadMusicFromDiskSubject() *pattern.OneArgSubject[[]Music] {
	return &state.onReadMusicsDiskSubject
}

func (state *State) OnSelectMusicSubject() *pattern.ThreeArgSubject[Album, []Music, Music] {
	return &state.onSelectMusicSubject
}

func (state *State) SetSelectedAlbum(album *Album) {
	state.onSelectAlbumSubject.NotifyAll(*album)
	if state.album != *album {
		state.album = *album
		state.musics = ReadMusicFromDisk(*album)
		state.onReadMusicsDiskSubject.NotifyAll(state.musics)
	}
}

func (state *State) SetSelectedMusic(music *Music) {
	state.onSelectMusicSubject.NotifyAll(state.album, state.musics, *music)
}

func ReadAlbumsFromDisk() []Album {
	directories, err := os.ReadDir(resource.GetAlbumRootPath())
	if err != nil {
		log.Fatal(err)
	}

	albums := []Album{}
	for _, directory := range directories {
		if directory.IsDir() {

			//read album config
			configPath := resource.GetAlbumConfigPath(directory.Name())
			config, err := os.ReadFile(configPath)
			if err != nil {
				log.Fatal(err)
			}

			//read last modified date
			info, err := os.Stat(configPath)
			if err != nil {
				log.Fatal(err)
			}

			musicNumber := bytes.Count(config, []byte{'\n'})
			albumCover := canvas.NewImageFromFile(resource.GetAlbumIconPath(directory.Name()))
			albums = append(albums, Album{directory.Name(), info.ModTime(), musicNumber, albumCover})
		}
	}
	return albums
}

func ReadMusicFromDisk(album Album) []Music {
	// a rough estimation of the music duration in nanoseconds
	estimateDuration := func(musicFileInfo fs.FileInfo) time.Duration {
		return time.Duration(musicFileInfo.Size() * MAGIC_RATIO / (AUDIO_BIT_DEPTH * NUM_OF_CHANNELS * SAMPLING_RATE))
	}

	config, err := os.ReadFile(resource.GetAlbumConfigPath(album.title))
	if err != nil {
		log.Fatal(err)
	}

	//read music name from config
	music := []Music{}
	scanner := bufio.NewScanner(bytes.NewReader(config))
	for scanner.Scan() {
		//open music file
		info, err := os.Stat(resource.GetMusicPath(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}

		music = append(music, Music{scanner.Text(), estimateDuration(info), info.ModTime()})
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	return music
}

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

	state.OnReadAlbumsFromDiskSubject().NotifyAll(ReadAlbumsFromDisk())
	return nil
}

func RenameAlbum(oldTitle, newTitle string) error {
	oldPath := resource.GetAlbumFolderPath(oldTitle)
	newPath := resource.GetAlbumFolderPath(newTitle)

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return fmt.Errorf("album \"%v\" is already existed", newTitle)
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	state.OnReadAlbumsFromDiskSubject().NotifyAll(ReadAlbumsFromDisk())
	return nil
}

func SetAlbumCover(title, coverIconPath string) error {
	coverIconData, err := os.ReadFile(coverIconPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(resource.GetAlbumIconPath(title), coverIconData, os.ModePerm); err != nil {
		return err
	}
	state.OnReadAlbumsFromDiskSubject().NotifyAll(ReadAlbumsFromDisk())
	return nil
}

func RemoveAlbum(title string) error {
	if err := os.RemoveAll(resource.GetAlbumFolderPath(title)); err != nil {
		return err
	}
	state.OnReadAlbumsFromDiskSubject().NotifyAll(ReadAlbumsFromDisk())
	return nil
}
