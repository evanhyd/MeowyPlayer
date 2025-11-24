package loggers

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	Log  *slog.Logger
	file *os.File
}

func MakeLogger() Logger {
	file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Panic(err)
	}
	handler := slog.NewJSONHandler(io.MultiWriter(os.Stderr, file), &slog.HandlerOptions{AddSource: true})
	return Logger{Log: slog.New(handler), file: file}
}

func (l *Logger) Close() error {
	return l.file.Close()
}
