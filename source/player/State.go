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

	"fyne.io/fyne/v2/canvas"
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/pattern"
	"meowyplayer.com/source/resource"
)

var state State

func init() {
	state = State{}
}

func GetState() *State {
	return &state
}

type State struct {
	album                       Album
	musics                      []Music
	onReadAlbumsFromDiskSubject pattern.OneArgSubject[[]Album]
	onSelectAlbumSubject        pattern.OneArgSubject[Album]
	onReadMusicsDiskSubject     pattern.OneArgSubject[[]Music]
	onSelectMusicSubject        pattern.ThreeArgSubject[Album, []Music, Music]
}

func (state *State) Album() Album {
	return state.album
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
		// var err error
		state.musics, _ = ReadMusicFromDisk(*album)
		state.onReadMusicsDiskSubject.NotifyAll(state.musics)
	}
}

func (state *State) SetSelectedMusic(music *Music) {
	state.onSelectMusicSubject.NotifyAll(state.album, state.musics, *music)
}

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
			albumCover := canvas.NewImageFromFile(resource.GetAlbumIconPath(directory.Name()))
			albumCover.SetMinSize(resource.GetAlbumCoverSize())
			albums = append(albums, Album{directory.Name(), info.ModTime(), musicNumber, albumCover})
		}
	}
	return albums, nil
}

func ReadMusicFromDisk(album Album) ([]Music, error) {
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
		//open music file
		info, err := os.Stat(resource.GetMusicPath(scanner.Text()))
		if err != nil {
			return nil, err
		}

		music = append(music, Music{scanner.Text(), estimateDuration(info), info.ModTime()})
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return music, nil
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

	albums, err := ReadAlbumsFromDisk()
	state.OnReadAlbumsFromDiskSubject().NotifyAll(albums)
	return err
}

func RenameAlbum(album Album, newTitle string) error {
	oldPath := resource.GetAlbumFolderPath(album.Title())
	newPath := resource.GetAlbumFolderPath(newTitle)

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return fmt.Errorf("album \"%v\" is already existed", newTitle)
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	albums, err := ReadAlbumsFromDisk()
	state.OnReadAlbumsFromDiskSubject().NotifyAll(albums)
	return err
}

func SetAlbumCover(album Album, coverIconPath string) error {
	coverIconData, err := os.ReadFile(coverIconPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(resource.GetAlbumIconPath(album.Title()), coverIconData, os.ModePerm); err != nil {
		return err
	}

	albums, err := ReadAlbumsFromDisk()
	state.OnReadAlbumsFromDiskSubject().NotifyAll(albums)
	return err
}

func RemoveAlbum(album Album) error {
	if err := os.RemoveAll(resource.GetAlbumFolderPath(album.Title())); err != nil {
		return err
	}

	albums, err := ReadAlbumsFromDisk()
	state.OnReadAlbumsFromDiskSubject().NotifyAll(albums)
	return err
}

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

	//update GUI
	albums, err := ReadAlbumsFromDisk()
	if err != nil {
		return err
	}
	index := slices.IndexFunc(albums, func(a Album) bool { return a.Title() == album.Title() })
	if index == -1 {
		return fmt.Errorf("can not find album %v", album.Title())
	}

	selectedAlbum := albums[index]
	state.OnReadAlbumsFromDiskSubject().NotifyAll(albums)
	state.SetSelectedAlbum(&selectedAlbum)
	return nil
}
