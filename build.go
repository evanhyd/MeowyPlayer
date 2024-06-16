package main

import (
	"flag"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var platform string
var runInDebug bool

func init() {
	flag.StringVar(&platform, "platform", runtime.GOOS, "The executable plaforms: windows, linux, darwin.")
	flag.BoolVar(&runInDebug, "debug", false, "Run the binary in debug mode after building.")
	flag.Parse()
}

func runAt(dir string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.String())
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func run(command string, args ...string) {
	runAt("", command, args...)
}

type bundleFormat struct {
	name string
	path string
}

func bundleResource(outputFile string, toBundles []bundleFormat) {
	fyneToolPath := filepath.Join(build.Default.GOPATH, "bin", "fyne")

	getAsset := func(asset string) string {
		return filepath.Join("asset", asset)
	}

	run(fyneToolPath, "bundle", "-o", outputFile, "--package", "resource", "--name", toBundles[0].name, getAsset(toBundles[0].path))

	for _, toBundle := range toBundles[1:] {
		run(fyneToolPath, "bundle", "-o", outputFile, "--append", "--name", toBundle.name, getAsset(toBundle.path))
	}
}

func buildBinary() {
	switch platform {
	case "windows":
		if runInDebug {
			runAt("source", "go", "build", "-o", filepath.Join("..", "meowyplayer.exe"), "main.go")
		} else {
			runAt("source", "go", "build", "-ldflags", "-H=windowsgui", "-o", filepath.Join("..", "meowyplayer.exe"), "main.go")
		}
	case "linux", "darwin":
		runAt("source", "go", "build", "-o", filepath.Join("..", "meowyplayer"), "main.go")
	default:
		panic("unknown platform")
	}

	if runInDebug {
		run("./meowyplayer")
	}
}

func main() {

	//bundle icons
	iconBundles := []bundleFormat{
		{"WindowIcon", "icon.png"}, //unfortunately fyne doesn't support svg as system tray icon
		{"CollectionTabIcon", "collection_tab.svg"},
		{"YouTubeIcon", "youtube.svg"},
		{"AlphabeticalIcon", "alphabetical.svg"},
	}
	iconPath := filepath.Join("source", "resource", "Icon.go")
	bundleResource(iconPath, iconBundles)

	//build or execute
	buildBinary()
}
