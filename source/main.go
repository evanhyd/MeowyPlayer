package main

import (
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2"
)

// Initialize a log file and redirect STDOUT, STDERR to it.
func initializeLogger() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fyne.LogError("failed to initialize log file", err)
		return
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))
}

func main() {
	initializeLogger()

}
