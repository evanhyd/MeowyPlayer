package main

import (
	"flag"
	"fmt"
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
			run("fyne", "package", "--release", "--src", "source", "--exe", kExeName) //-o has missing icon bug
			os.Rename(filepath.Join("source", kExeName), filepath.Join(".", kExeName))
		} else {
			runAt("source", "go", "build", "-o", filepath.Join("..", kExeName), "main.go")
			run("./meowyplayer")
		}
	} else {
		fmt.Println("currently does not support other platform")
	}
}
