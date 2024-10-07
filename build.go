package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var releaseFlag bool
var platformFlag string

func init() {
	flag.BoolVar(&releaseFlag, "release", false, "Compile the build in release mode.")
	flag.StringVar(&platformFlag, "platform", "desktop", "Target platform (desktop, android).")
	flag.Parse()
}

func runAt(dir string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	log.Println(cmd.String())
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func run(command string, args ...string) {
	runAt("", command, args...)
}

func main() {
	const kExeName = "meowyplayer.exe"

	if platformFlag == "desktop" {
		if releaseFlag {
			run("fyne", "package", "--src", "source", "--exe", kExeName, "--release") //-o has missing icon bug
			run("mv", filepath.Join("source", kExeName), ".")
		} else {
			runAt("source", "go", "build", "-o", filepath.Join("..", kExeName), "main.go")
			run("./meowyplayer")
		}
	} else {
		runAt("source", "fyne", "package", "--os", platformFlag, "--release")
		run("mv", filepath.Join("source", kExeName), ".")
	}
}
