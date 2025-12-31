package ui

import (
	"embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

//go:embed internal/assets/translations
var translations embed.FS

func RunApp() {
	lang.AddTranslationsFS(translations, ".")

	mainApp := app.New()
	mainApp.SetIcon(resourceIconPng)
	mainApp.Settings().SetTheme(newVanillaTheme())

	mainWindow := mainApp.NewWindow(lang.L("MeowyPlayer"))
	mainWindow.Resize(fyne.NewSize(809, 500))

	appTab := container.NewAppTabs(container.NewTabItemWithIcon(lang.L("Explore"), theme.MediaMusicIcon(), newExplorePage()))
	appTab.SetTabLocation(container.TabLocationLeading)
	mainWindow.SetContent(appTab)

	// System tray.
	// mainWindow.SetCloseIntercept(mainWindow.Hide)
	// if desktop, ok := mainApp.(desktop.App); ok {
	// 	desktop.SetSystemTrayMenu(fyne.NewMenu("",
	// 		fyne.NewMenuItem("Show", mainWindow.Show),
	// 	))
	// }
	mainWindow.ShowAndRun()
}
