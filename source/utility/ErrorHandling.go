package utility

import (
	"flag"
	"io"
	"log"
	"os"
)

func InitLogger() {
	const kLogFileName = "log.txt"
	stdlog := flag.Bool("stdlog", true, "output logs to standard error stream")
	filelog := flag.Bool("filelog", true, "output logs to "+kLogFileName)
	flag.Parse()

	//set up writing destination
	loggers := []io.Writer{}

	if *stdlog {
		loggers = append(loggers, os.Stderr)
	}

	if *filelog {
		file, err := os.OpenFile(kLogFileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
		MustOk(err)
		loggers = append(loggers, file)
	}

	log.SetOutput(io.MultiWriter(loggers...))
}

// Error must not occur, for pre/post conditions checking
func MustOk(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// Error may occur, ex: IO fails due to unavoidable network error
func ShouldOk(err error) {
	if err != nil {
		log.Println(err)
	}
}

func MustNotNil(object any) {
	if object == nil {
		log.Panicf("%v is nil\n", object)
	}
}
