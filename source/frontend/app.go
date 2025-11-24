package frontend

import (
	"embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
)

//go:embed internal/assets/translations
var translations embed.FS

func RunApp() {
	lang.AddTranslationsFS(translations, ".")

	mainApp := app.New()
	mainApp.SetIcon(resourceIconPng)
	mainApp.Settings().SetTheme(newVanillaTheme())

	mainWindow := mainApp.NewWindow(lang.L("MeowyPlayer"))
	// mainWindow.SetFullScreen(true)
	mainWindow.SetContent(canvas.NewRectangle(color.White)) //TODO:

	// System tray.
	mainWindow.SetCloseIntercept(mainWindow.Hide)
	if desktop, ok := mainApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("",
			fyne.NewMenuItem("Show", mainWindow.Show),
		))
	}

	mainWindow.ShowAndRun()
}
