package logger

import (
	"flag"
	"fmt"
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
			Error(fmt.Errorf("%v %v", err.Error(), "failed to initiate logger"), 0)
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

func Error(err error, stackRewind int) {
	_, file, line, ok := runtime.Caller(stackRewind + 1)
	if ok {
		log.Printf("%s[%d] - %v\n", file, line, err)
	}
}

func Fatal(err error, stackRewind int) {
	_, file, line, ok := runtime.Caller(stackRewind + 1)
	if ok {
		log.Panicf("%s[%d] - %v\n", file, line, err)
	}
}
