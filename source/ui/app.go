package ui

import (
	"embed"
	"meowyplayer/context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

//go:embed internal/translations
var translations embed.FS

func RunApp(userContext *context.UserContext, postUICallback func()) {
	lang.AddTranslationsFS(translations, ".")

	mainApp := app.New()
	mainApp.SetIcon(resourceIconPng)
	mainApp.Settings().SetTheme(newVanillaTheme())

	mainWindow := mainApp.NewWindow(lang.L("MeowyPlayer"))
	mainWindow.Resize(fyne.NewSize(809, 500))

	appTab := container.NewAppTabs(
		container.NewTabItemWithIcon(lang.L("Explore"), theme.MediaMusicIcon(), newExplorePage(userContext)),
		container.NewTabItemWithIcon(lang.L("Playlist"), resourcePlaylistSvg, newPlaylistPage(userContext)),
	)
	appTab.SetTabLocation(container.TabLocationLeading)
	mainWindow.SetContent(appTab)

	// System tray.
	// mainWindow.SetCloseIntercept(mainWindow.Hide)
	// if desktop, ok := mainApp.(desktop.App); ok {
	// 	desktop.SetSystemTrayMenu(fyne.NewMenu("",
	// 		fyne.NewMenuItem("Show", mainWindow.Show),
	// 	))
	// }

	postUICallback()
	mainWindow.ShowAndRun()
}
