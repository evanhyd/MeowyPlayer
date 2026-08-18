package ui

import (
	"embed"
	"meowyplayer/mcontext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

//go:embed internal/translations
var translations embed.FS

func RunApp(userContext *mcontext.UserContext) {
	// Localization.
	lang.AddTranslationsFS(translations, ".")

	// Main app.
	meowApp := app.New()
	meowApp.SetIcon(resourceIconPng)
	meowApp.Settings().SetTheme(newVanillaTheme())

	// Auth window.
	authWindow := newAuthWindow(meowApp, userContext)

	// Switch to the music player window if login successfully.
	userContext.AddListener(mcontext.OnPutUserEvent, func(any) {
		musicWin := newMusicWindow(meowApp, userContext)
		authWindow.Close()
		musicWin.Show()
	})

	authWindow.ShowAndRun()
}

func newAuthWindow(meowApp fyne.App, userContext *mcontext.UserContext) fyne.Window {
	win := meowApp.NewWindow(lang.L("MeowyPlayer Auth Window"))
	win.Resize(fyne.NewSize(550, 250))
	win.SetContent(newAuthPage(userContext))
	win.CenterOnScreen()
	win.SetFixedSize(true)
	return win
}

func newMusicWindow(meowApp fyne.App, userContext *mcontext.UserContext) fyne.Window {
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
