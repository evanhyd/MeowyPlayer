package logger

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime"
)

func Initiate() {
	const kLogFileName = "log.txt"
	errLog := flag.Bool("errlog", true, "output logs to the error stream")
	fileLog := flag.Bool("filelog", true, "output logs to "+kLogFileName)
	flag.Parse()

	//set up writing destination
	loggers := []io.Writer{}

	if *errLog {
		loggers = append(loggers, os.Stderr)
	}

	if *fileLog {
		file, err := os.OpenFile(kLogFileName, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			Error("failed to initiate logger", err, 1)
		}
		loggers = append(loggers, file)
	}

	log.SetOutput(io.MultiWriter(loggers...))
}

func Error(message string, err error, frameRewind int) {
	_, file, line, ok := runtime.Caller(frameRewind)
	if ok {
		log.Printf("error occured at: %s:%d\n%v: %v\n", file, line, message, err)
	}
}
