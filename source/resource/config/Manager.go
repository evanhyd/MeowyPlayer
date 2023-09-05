package config

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/resource/path"
	"meowyplayer.com/source/resource/texture"
	"meowyplayer.com/source/utility"
)

var currentConfig player.Config

func Current() *player.Config {
	return &currentConfig
}

func Set(config *player.Config) {
	currentConfig.Date = config.Date
	currentConfig.Albums = config.Albums
	currentConfig.NotifyAll()
}

func Reload() error {
	if err := SaveToLocal(&currentConfig); err != nil {
		return err
	}

	config, err := LoadFromLocal()
	if err != nil {
		return err
	}

	Set(&config)
	return err
}

func LoadFromLocal() (player.Config, error) {
	config := player.Config{}
	if err := utility.ReadJson(path.Config(), &config); err != nil {
		return config, err
	}

	//load icons
	getIcon := func(album *player.Album) fyne.Resource {
		const missingTexturePath = "missing_texture.png"

		//if fail, then load the placeholder texture
		icon, err := fyne.LoadResourceFromPath(path.Icon(album))
		if os.IsNotExist(err) {
			return texture.Get(missingTexturePath)
		}
		utility.MustOk(err)
		return icon
	}

	for i := range config.Albums {
		config.Albums[i].Cover = canvas.NewImageFromResource(getIcon(&config.Albums[i]))
	}

	config.NotifyAll()
	return config, nil
}

func SaveToLocal(config *player.Config) error {
	return utility.WriteJson(path.Config(), config)
}

func Download() player.Config {
	//send request to the server
	panic("not implemented")
}

func Upload() {
	//upload to the server
	panic("not implemented")
}
