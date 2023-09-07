package manager

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"golang.org/x/exp/slices"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource/path"
	"meowyplayer.com/source/resource/texture"
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

	//load icons
	getCover := func(album *player.Album) fyne.Resource {
		const missingTexturePath = "missing_texture.png"

		//if fail, then load the placeholder texture
		icon, err := fyne.LoadResourceFromPath(path.Icon(album))
		if os.IsNotExist(err) {
			return texture.Get(missingTexturePath)
		}
		utility.MustOk(err)
		return icon
	}

	for i := range inUse.Albums {
		inUse.Albums[i].Cover = canvas.NewImageFromResource(getCover(&inUse.Albums[i]))
	}

	return inUse, nil
}

func SaveToLocalConfig(config *player.Config) error {
	return utility.WriteJson(path.Config(), config)
}

func ReloadConfig() error {
	if err := SaveToLocalConfig(configData.Get()); err != nil {
		return err
	}

	config, err := LoadFromLocalConfig()
	if err != nil {
		return err
	}

	configData.Set(&config)
	return err
}

func ReloadAlbum() error {
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

	file, err := os.Create(path.Icon(&album))
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, iconImage); err != nil {
		return err
	}

	return ReloadConfig()
}

func DeleteAlbum(album *player.Album) error {
	source := getSourceAlbum(album)
	inUse := configData.Get()
	index := slices.IndexFunc(inUse.Albums, func(a player.Album) bool { return a.Title == source.Title })
	last := len(inUse.Albums) - 1

	//remove album icon
	if err := os.Remove(path.Icon(source)); err != nil && !os.IsNotExist(err) {
		return err
	}

	inUse.Albums[index] = inUse.Albums[last]
	inUse.Albums = inUse.Albums[:last]

	return ReloadConfig()
}

func DeleteMusic(music *player.Music) error {
	source := getSourceAlbum(albumData.Get())
	index := slices.IndexFunc(source.MusicList, func(m player.Music) bool { return m.Title == music.Title })
	last := len(source.MusicList) - 1

	source.MusicList[index] = source.MusicList[last]
	source.MusicList = source.MusicList[:last]

	if err := ReloadConfig(); err != nil {
		return err
	}

	return ReloadAlbum()
}

func UpdateTitle(album *player.Album, title string) error {
	inUse := configData.Get()
	if slices.ContainsFunc(inUse.Albums, func(a player.Album) bool { return a.Title == title }) {
		return fmt.Errorf("album \"%v\" already exists", title)
	}

	source := getSourceAlbum(album)
	source.Date = time.Now()
	inUse.Date = time.Now()

	oldPath := path.Icon(source)
	source.Title = title
	if err := os.Rename(oldPath, path.Icon(source)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return ReloadConfig()
}

func UpdateCover(album *player.Album, iconPath string) error {
	source := getSourceAlbum(album)
	source.Date = time.Now() //update description -> update album view
	configData.Get().Date = time.Now()

	icon, err := os.ReadFile(iconPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Icon(source), icon, os.ModePerm)
	if err != nil {
		return err
	}

	return ReloadConfig()
}
