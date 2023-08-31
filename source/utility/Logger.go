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
