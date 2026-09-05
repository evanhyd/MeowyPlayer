package ui

import (
	"embed"
	"log/slog"
	"meowyplayer/storages"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

//go:embed internal/translations
var translations embed.FS

func RunApp(userContext *storages.UserContext) {
	// Localization.
	err := lang.AddTranslationsFS(translations, path.Join("internal", "translations"))
	if err != nil {
		slog.Error("failed to add translation", "error", err)
	}

	// Main app.
	meowApp := app.New()
	meowApp.SetIcon(resourceIconPng)
	meowApp.Settings().SetTheme(newVanillaTheme())

	// Check if already login. This bypasses the server check and allows the user to access the music player offline.
	if _, err := userContext.GetUser(); err != nil {
		authWindow := newAuthWindow(meowApp, userContext)
		userContext.AddListener(storages.OnPutUserEvent, func(any) {
			musicWin := newMusicWindow(meowApp, userContext)
			authWindow.Close()
			musicWin.Show()
			userContext.OnInit()
		})
		authWindow.ShowAndRun()

	} else {
		musicWin := newMusicWindow(meowApp, userContext)
		musicWin.Show()
		userContext.OnInit()
		meowApp.Run()
	}
}

func newAuthWindow(meowApp fyne.App, userContext *storages.UserContext) fyne.Window {
	win := meowApp.NewWindow(lang.L("MeowyPlayer Auth Window"))
	win.Resize(fyne.NewSize(750, 250))
	win.SetContent(newAuthPage(userContext))
	win.CenterOnScreen()
	return win
}

func newMusicWindow(meowApp fyne.App, userContext *storages.UserContext) fyne.Window {
	win := meowApp.NewWindow(lang.L("MeowyPlayer"))
	win.SetMaster()
	win.Resize(fyne.NewSize(800, 550))

	// Content.
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(lang.L("Explore"), theme.MediaMusicIcon(), newExplorePage(userContext)),
		container.NewTabItemWithIcon(lang.L("Playlist"), resourcePlaylistSvg, container.NewStack(newPlaylistPage(userContext), newMusicPage(userContext))),
		container.NewTabItemWithIcon(lang.L("Profile"), theme.AccountIcon(), newProfilePage(userContext)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	tabs.SelectIndex(1)
	controller := newMusicController(userContext)
	win.SetContent(container.NewBorder(nil, controller, nil, nil, tabs))

	// System tray menu.
	if desktop, ok := meowApp.(desktop.App); ok {
		desktop.SetSystemTrayMenu(fyne.NewMenu("", fyne.NewMenuItem("Show", win.Show)))
	}

	return win
}
