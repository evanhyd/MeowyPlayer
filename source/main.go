package main

import (
	"io"
	"log"
	"os"
)

// Initialize a log file and redirect STDOUT, STDERR to it.
func initializeLogger() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("failed to initialize log file", err)
		return
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))
}

func main() {
	initializeLogger()
}
