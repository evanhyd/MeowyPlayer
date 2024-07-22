package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var platform string
var release bool

func init() {
	flag.StringVar(&platform, "platform", runtime.GOOS, "The executable plaforms: windows, linux, darwin.")
	flag.BoolVar(&release, "debug", false, "Run the binary in debug mode after building.")
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
	if release {
		run("fyne", "package", "--src", "source", "--exe", "..", "--release")
	} else {
		run("fyne", "package", "--src", "source", "--exe", "..")
	}
}
