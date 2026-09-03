package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func runAt(path string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(command string, args ...string) {
	runAt(".", command, args...)
}

func compileWindows(isRelease bool) {
	if isRelease {
		run("fyne", "package", "--src", "source", "--exe", "..", "--release")
	} else {
		run("fyne", "package", "--src", "source", "--exe", "..")
	}
}

func compileAndroid(isRelease bool) {
	if isRelease {
		runAt("./source", "fyne", "package", "-os", "android", "-release")
	} else {
		runAt("./source", "fyne", "package", "-os", "android")
	}
	os.Rename("./source/meowyplayer.apk", "meowyplayer.apk")
}

func main() {
	os := flag.String("os", "windows", "Target platform operating system.")
	isRelease := flag.Bool("release", false, "Build in release mode. This optimizes the app but drops lots of debugging symbols.")
	flag.Parse()

	// Package assets.
	runAt("./source/ui", "fyne", "bundle", "-o", "assets.go", "--pkg", "ui", "internal/assets/*")

	switch *os {
	case "windows":
		compileWindows(*isRelease)
	case "android":
		compileAndroid(*isRelease)
	default:
	}

	// // Force Linux slashes for Docker
	// assetsPath := "./source/ui/assets.go"
	// if data, err := os.ReadFile(assetsPath); err == nil {
	// 	fixedContent := strings.ReplaceAll(string(data), "\\", "/")
	// 	os.WriteFile(assetsPath, []byte(fixedContent), 0644)
	// }

	// // Build.
	// args := []string{*platform}
	// if *release {
	// 	args = append(args, "-release")
	// }
	// args = append(args)
	// runAt("./source", "fyne-cross", args...)
}
