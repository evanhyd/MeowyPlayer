package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var releaseFlag bool

func init() {
	flag.BoolVar(&releaseFlag, "release", false, "Compile the build in release mode.")
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
	if releaseFlag {
		// run("fyne", "package", "--src", "source", "--exe", filepath.Join("..", kExeName), "--release") //missing icon bug
		run("fyne", "package", "--src", "source", "--exe", kExeName, "--release")
		run("mv", filepath.Join("source", kExeName), ".")
	} else {
		runAt("source", "go", "build", "-o", filepath.Join("..", kExeName), "main.go")
		run("./meowyplayer")
	}
}
