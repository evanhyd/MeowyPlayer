package ui

import (
	"embed"
	"meowyplayer/context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
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
	mainWindow.SetCloseIntercept(mainApp.Quit)
	if desktop, ok := mainApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", mainWindow.Show)))
	}

	appTab := container.NewAppTabs(
		container.NewTabItemWithIcon(lang.L("Explore"), theme.MediaMusicIcon(), newExplorePage(userContext)),
		container.NewTabItemWithIcon(lang.L("Playlist"), resourcePlaylistSvg, container.NewStack(newPlaylistPage(userContext), newMusicPage(userContext))),
	)
	appTab.SetTabLocation(container.TabLocationLeading)
	appTab.SelectIndex(1)

	postUICallback()
	mainWindow.SetContent(container.NewBorder(nil, newMusicController(userContext), nil, nil, appTab))
	mainWindow.ShowAndRun()
}
