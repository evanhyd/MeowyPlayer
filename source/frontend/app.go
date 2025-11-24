package frontend

import (
	"embed"

	"fyne.io/fyne/v2/lang"
)

//go:embed internal/assets/translations
var assets embed.FS

func RunApp() {
	lang.AddTranslationsFS(assets, ".")

}
