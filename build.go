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

func main() {
	release := flag.Bool("release", false, "Build in release mode. This optimizes the app but drops lots of debugging symbols.")
	flag.Parse()

	runAt("./source/ui", "fyne", "bundle", "-o", "assets.go", "--pkg", "ui", "internal/assets/*")
	if *release {
		run("fyne", "package", "--src", "source", "--exe", "..", "--release")
	} else {
		run("fyne", "package", "--src", "source", "--exe", "..")
	}
}
