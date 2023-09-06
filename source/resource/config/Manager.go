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

var configData utility.Data[player.Config]

func Get() *utility.Data[player.Config] {
	return &configData
}

func Reload() error {
	if err := SaveToLocal(configData.Get()); err != nil {
		return err
	}

	config, err := LoadFromLocal()
	if err != nil {
		return err
	}

	configData.Set(&config)
	return err
}

func LoadFromLocal() (player.Config, error) {
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
