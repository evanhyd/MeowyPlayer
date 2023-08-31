package resource

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
)

const (
	albumPath          = "album"
	iconPath           = "icon"
	configFile         = "config.json"
	texturePath        = "texture"
	missingTexturePath = "missing_texture.png"
)

func GetConfigPath() string {
	return filepath.Join(albumPath, configFile)
}

func GetIconPath(album *player.Album) string {
	return filepath.Join(albumPath, iconPath, album.Title+".png")
}

func GetIcon(album *player.Album) (fyne.Resource, error) {
	return fyne.LoadResourceFromPath(GetIconPath(album))
}

func SetIcon(album *player.Album, iconPath string) error {
	icon, err := os.ReadFile(iconPath)
	if err != nil {
		return err
	}
	return os.WriteFile(GetIconPath(album), icon, os.ModePerm)
}

func GetTexture(textureName string) fyne.Resource {
	texture, err := fyne.LoadResourceFromPath(filepath.Join(texturePath, textureName))

	//if fail, then load the placeholder texture
	if os.IsNotExist(err) {
		texture, err = fyne.LoadResourceFromPath(filepath.Join(texturePath, missingTexturePath))
	}

	utility.MustOk(err)
	return texture
}
