package utility

import (
	"flag"
	"io"
	"log"
	"os"
)

func InitLogger() {
	const kLogFileName = "log.txt"
	stdLog := flag.Bool("stdlog", true, "output logs to the standard error stream")
	fileLog := flag.Bool("filelog", true, "output logs to "+kLogFileName)
	flag.Parse()

	//set up writing destination
	loggers := []io.Writer{}

	if *stdLog {
		loggers = append(loggers, os.Stderr)
	}

	if *fileLog {
		file, err := os.OpenFile(kLogFileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
		MustNil(err)
		loggers = append(loggers, file)
	}
	log.SetOutput(io.MultiWriter(loggers...))
}

// Must satisfy, used as pre/post conditions check.
func MustNil(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func MustNotNil(object any) {
	if object == nil {
		log.Panicf("%v is nil\n", object)
	}
}

// Should satisfy, however, it is legal that the program fails to satisfy such condition.
func ShouldNil(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ShouldNotNil(err error) {
	if err == nil {
		log.Println(err)
	}
}

// Program invariant assertion, must be true
func Assert(condition func() bool) {
	if !condition() {
		log.Panic("condition failed\n")
	}
}
