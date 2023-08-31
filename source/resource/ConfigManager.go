package resource

import (
	"os"

	"fyne.io/fyne/v2/canvas"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
)

var localConfig player.Config

func GetCurrentConfig() *player.Config {
	return &localConfig
}

func SetCurrentConfig(config *player.Config) {
	localConfig.Date = config.Date
	localConfig.Albums = config.Albums
	localConfig.NotifyAll()
}

func ReloadCurrentConfig() error {
	if err := SaveToLocalConfig(&localConfig); err != nil {
		return err
	}

	config, err := LoadFromLocalConfig()
	if err == nil {
		SetCurrentConfig(&config)
	}
	return err
}

func LoadFromLocalConfig() (player.Config, error) {
	config := player.Config{}
	if err := utility.ReadJson(GetConfigPath(), &config); err != nil {
		return config, nil
	}

	//load icons
	for i := range config.Albums {
		resource, err := GetIcon(&config.Albums[i])
		if err != nil {
			if os.IsNotExist(err) {
				resource = GetTexture(missingTexturePath)
			} else {
				return config, err
			}
		}
		config.Albums[i].Cover = canvas.NewImageFromResource(resource)
	}
	config.NotifyAll()
	return config, nil
}

func SaveToLocalConfig(config *player.Config) error {
	return utility.WriteJson(GetConfigPath(), config)
}

func DownloadConfig() player.Config {
	//send request to the server
	panic("not implemented")
}

func UploadConfig() {
	//upload to the server
	panic("not implemented")
}
