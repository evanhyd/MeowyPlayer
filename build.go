package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func run(dir string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	log.Println(cmd.String())

	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	var releaseFlag bool
	flag.BoolVar(&releaseFlag, "release", false, "Compile the build in release mode.")
	flag.Parse()

	const outputName = "meowyplayer.exe"

	if releaseFlag {
		run("", "fyne", "package", "--release", "--src", "source", "--exe", outputName) //-o has missing icon bug
		os.Rename(filepath.Join("source", outputName), filepath.Join(".", outputName))
	} else {
		run("source", "go", "build", "-o", filepath.Join("..", outputName), "main.go")
		run("", "./meowyplayer")
	}
}
