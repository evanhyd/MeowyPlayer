package manager

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/path"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource"
	"meowyplayer.com/source/utility"
)

var configData utility.Data[player.Config]
var albumData utility.Data[player.Album]

// the album pointer parameter may refer to a temporary object from the view list
// we need the original one from the config
func getSourceAlbum(album *player.Album) *player.Album {
	inUse := configData.Get()
	index := slices.IndexFunc(inUse.Albums, func(a player.Album) bool { return a.Title == album.Title })
	return &inUse.Albums[index]
}

func GetCurrentConfig() *utility.Data[player.Config] {
	return &configData
}

func GetCurrentAlbum() *utility.Data[player.Album] {
	return &albumData
}

func LoadFromLocalConfig() (player.Config, error) {
	inUse := player.Config{}
	if err := utility.ReadJson(path.Config(), &inUse); err != nil {
		return inUse, err
	}

	for i := range inUse.Albums {
		album := &inUse.Albums[i]

		//read icon
		album.Cover = resource.GetCover(album)

		//read music file size
		for j := range album.MusicList {
			music := &album.MusicList[j]
			fileInfo, err := os.Stat(path.Music(music))
			utility.ShouldNil(err)
			music.FileSize = fileInfo.Size()
		}
	}

	return inUse, nil
}

func reloadConfig() error {
	if err := utility.WriteJson(path.Config(), configData.Get()); err != nil {
		return err
	}

	config, err := LoadFromLocalConfig()
	if err != nil {
		return err
	}

	configData.Set(&config)
	return nil
}

func reloadAlbum() error {
	source := getSourceAlbum(albumData.Get())
	albumData.Set(source)
	return nil
}

func AddAlbum() error {
	inUse := configData.Get()

	//generate title
	title := ""
	for i := 0; i < math.MaxInt; i++ {
		title = fmt.Sprintf("My Album (%v)", i)
		if !slices.ContainsFunc(inUse.Albums, func(a player.Album) bool { return a.Title == title }) {
			break
		}
	}

	//generate album
	album := player.Album{Date: time.Now(), Title: title}
	inUse.Albums = append(inUse.Albums, album)

	//generate album cover
	iconColor := color.NRGBA{uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32()), uint8(rand.Uint32())}
	iconImage := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	iconImage.SetNRGBA(0, 0, iconColor)

	imageData := bytes.Buffer{}
	if err := png.Encode(&imageData, iconImage); err != nil {
		return err
	}

	if err := os.WriteFile(path.Cover(&album), imageData.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return reloadConfig()
}

func AddMusic(musicInfo fyne.URIReadCloser) error {
	music := player.Music{Date: time.Now(), Title: musicInfo.URI().Name()}
	source := getSourceAlbum(albumData.Get())
	source.MusicList = append(source.MusicList, music)

	//copy the file to the music directory
	musicFile, err := os.ReadFile(musicInfo.URI().Path())
	if err != nil {
		return err
	}

	if err = os.WriteFile(path.Music(&music), musicFile, os.ModePerm); err != nil {
		return err
	}

	if err := reloadConfig(); err != nil {
		return err
	}
	return reloadAlbum()
}

func DeleteAlbum(album *player.Album) error {
	source := getSourceAlbum(album)
	inUse := configData.Get()
	index := slices.IndexFunc(inUse.Albums, func(a player.Album) bool { return a.Title == source.Title })
	last := len(inUse.Albums) - 1

	//remove album icon
	if err := os.Remove(path.Cover(source)); err != nil && !os.IsNotExist(err) {
		return err
	}

	inUse.Albums[index] = inUse.Albums[last]
	inUse.Albums = inUse.Albums[:last]

	return reloadConfig()
}

func DeleteMusic(music *player.Music) error {
	source := getSourceAlbum(albumData.Get())
	index := slices.IndexFunc(source.MusicList, func(m player.Music) bool { return m.Title == music.Title })
	last := len(source.MusicList) - 1

	source.MusicList[index] = source.MusicList[last]
	source.MusicList = source.MusicList[:last]

	if err := reloadConfig(); err != nil {
		return err
	}
	return reloadAlbum()
}

func UpdateTitle(album *player.Album, title string) error {
	inUse := configData.Get()
	if slices.ContainsFunc(inUse.Albums, func(a player.Album) bool { return a.Title == title }) {
		return fmt.Errorf("album \"%v\" already exists", title)
	}

	source := getSourceAlbum(album)
	source.Date = time.Now()
	inUse.Date = time.Now()

	oldPath := path.Cover(source)
	source.Title = title
	if err := os.Rename(oldPath, path.Cover(source)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return reloadConfig()
}

func UpdateCover(album *player.Album, iconPath string) error {
	source := getSourceAlbum(album)
	source.Date = time.Now() //update description -> update album view
	configData.Get().Date = time.Now()

	icon, err := os.ReadFile(iconPath)
	if err != nil {
		return err
	}

	if err = os.WriteFile(path.Cover(source), icon, os.ModePerm); err != nil {
		return err
	}

	return reloadConfig()
}
