package texture

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"meowyplayer.com/source/utility"
)

const (
	texturePath = "texture"
)

func Get(textureName string) fyne.Resource {
	texture, err := fyne.LoadResourceFromPath(filepath.Join(texturePath, textureName))
	utility.MustOk(err)
	return texture
}
