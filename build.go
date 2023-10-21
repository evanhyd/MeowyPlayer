package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func runCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	println(cmd.String())
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func runCmdWithDir(dir string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	println(cmd.String())
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	fynePath := filepath.Join(os.Getenv("GOPATH"), "bin", "fyne")

	bundleResource := func(file string, name string, resource string) {
		runCmd(fynePath, "bundle", "-o", file, "--package", "resource", "--name", name, resource)
	}

	appendResource := func(file string, name string, resource string) {
		runCmd(fynePath, "bundle", "-o", file, "--append", "--name", name, resource)
	}

	iconFile := filepath.Join("source", "core", "resource", "Icon.go")
	bundleResource(iconFile, "MissingIcon", filepath.Join("asset", "missing_asset.png"))
	appendResource(iconFile, "WindowIcon", filepath.Join("asset", "icon.ico"))
	appendResource(iconFile, "AlbumTabIcon", filepath.Join("asset", "album_tab.png"))
	appendResource(iconFile, "AlbumAdderOnlineIcon", filepath.Join("asset", "album_adder_online.png"))
	appendResource(iconFile, "MusicTabIcon", filepath.Join("asset", "music_tab.png"))
	appendResource(iconFile, "MusicAdderOnlineIcon", filepath.Join("asset", "music_adder_online.png"))
	appendResource(iconFile, "DefaultIcon", filepath.Join("asset", "default.png"))
	appendResource(iconFile, "RandomIcon", filepath.Join("asset", "random.png"))
	appendResource(iconFile, "YouTubeIcon", filepath.Join("asset", "youtube.png"))
	appendResource(iconFile, "BiliBiliIcon", filepath.Join("asset", "bilibili.png"))

	fontFile := filepath.Join("source", "core", "resource", "Font.go")
	bundleResource(fontFile, "RegularFont", filepath.Join("asset", "regular_font.ttf"))
	appendResource(fontFile, "BoldFont", filepath.Join("asset", "bold_font.ttf"))
	appendResource(fontFile, "ItalicFont", filepath.Join("asset", "italic_font.ttf"))
	appendResource(fontFile, "BoldItalicFont", filepath.Join("asset", "bold_italic_font.ttf"))

	// build binary
	platform := runtime.GOOS
	if len(os.Args) >= 2 {
		platform = os.Args[1]
	}

	switch platform {
	case "windows":
		runCmdWithDir("source", "go", "build", "-ldflags", "-H=windowsgui", "-o", "meowyplayer.exe", "main.go")
	case "linux", "darwin":
		runCmdWithDir("source", "go", "build", "-o", "meowyplayer", "main.go")
	default:
		panic("unknown platform")
	}
}
