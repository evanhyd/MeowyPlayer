package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func run(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	release := flag.Bool("release", false, "Build in release mode. This optimizes the app but drops lots of debugging symbols.")
	flag.Parse()

	if *release {
		run("fyne", "package", "--src", "source", "--exe", "..", "--release")
	} else {
		run("fyne", "package", "--src", "source", "--exe", "..")
		run("./meowyplayer")
	}
}
