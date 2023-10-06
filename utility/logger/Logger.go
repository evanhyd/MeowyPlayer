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
	fileLog := flag.Bool("filelog", true, "output logs to "+kLogFileName)
	errLog := flag.Bool("errlog", true, "output logs to the error stream")
	flag.Parse()

	//set up writing destination
	loggers := []io.Writer{}

	if *fileLog {
		file, err := os.OpenFile(kLogFileName, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			Error(err, "failed to initiate logger", 1)
		}
		loggers = append(loggers, file)
	}

	if *errLog {
		loggers = append(loggers, os.Stderr)
	}

	//if -H=windowsgui is set, then writing to Stderr stream will fails AND
	//stops all the subsequence writing operations to the remaining writers
	//hence log file should be written before the console output
	log.SetOutput(io.MultiWriter(loggers...))
}

func Error(err error, message string, frameRewind int) {
	_, file, line, ok := runtime.Caller(frameRewind)
	if ok {
		log.Printf("error occured at: %s:%d\n%v: %v\n", file, line, message, err)
	}
}
